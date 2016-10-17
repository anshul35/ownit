package Auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	//	"fmt"

	"github.com/anshul35/ownit/Auth/JWT"
	"github.com/anshul35/ownit/Models"
	"github.com/anshul35/ownit/Settings"

	log "github.com/golang/glog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     Settings.Config.Facebook["clientID"],
		ClientSecret: Settings.Config.Facebook["clientSecret"],
		RedirectURL:  Settings.Config.Facebook["redirectURL"],
		Scopes:       []string{"public_profile", "user_friends", "email", "user_birthday", "user_location"},
		Endpoint:     facebook.Endpoint,
	}
	oauthStateString = "thisshouldberandom"
)

func HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Fabook Login: Parse Facebook API URL Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Problem has been reported. Please try after some time!"))
		return
	}
	parameters := url.Values{}
	parameters.Add("app_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Error("Facebook Login: Invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Error("Facebook Login oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		log.Error("Facebook Login Get User using AccessToken: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	type Data struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}
	data := Data{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Error("Facebook Login Unable to decode User profile data. Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please try again after some time"))
		return
	}

	log.Info("Facebook Login Succesfull for user: " + data.Name)
	user := Models.User{UserID: data.Id, Name: data.Name}
	jwtToken, err := JWT.GetJWTToken(&user)
	defer user.Save()

	type RData struct {
		Token  string
		UserID string
	}
	rdata := RData{Token: jwtToken, UserID: user.UserID}
	err = json.NewEncoder(w).Encode(rdata)
	if err != nil {
		log.Error("Facebook Login: Error in encoding json. Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please try again after some time"))
		return
	}
	return
}
