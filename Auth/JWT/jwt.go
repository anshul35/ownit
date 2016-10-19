package JWT

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/golang/glog"

	"github.com/anshul35/ownit/Models"
	"github.com/anshul35/ownit/Settings/Constants"
	"github.com/anshul35/ownit/Utilities"
)

func GetJWTToken(user *Models.User) (string, error) {
	duration, err := time.ParseDuration(Constants.TokenExpiryDuration)
	if err != nil {
		log.Error("JWT Token: Unable to parse expiry duration constant. Error: ", err)
		return "", err
	}
	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		jwt.MapClaims{
			"userid": user.UserID,
			"exp":    time.Now().Add(duration).Unix(),
			"iat":    time.Now().Unix(),
		})
	key := Constants.ClientTokenSecret
	tokenString, err := token.SignedString([]byte(key))
	user.JWTToken = tokenString
	if err != nil {
		log.Error("JWT Token: Unable to sign the token with claims. Error: ", err)
	}
	return tokenString, err
}

func AuthenticateClientRequest(r *http.Request) error {
	params := r.URL.Query()
	p, ok := params["jwt_token"]
	if !ok {
		return errors.New("Need an access token for authentication!")
	}
	tokenString := p[0]
	token, err := jwt.Parse(tokenString, func(tok *jwt.Token) (interface{}, error) {
		secret := Constants.ClientTokenSecret
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Token verified
		fmt.Println(claims)
		userID := claims["userid"]
		user, err := Models.GetUserByID(userID.(string))
		newToken, err := GetJWTToken(user)
		user.JWTToken = newToken
		if err != nil {
			log.Error("Authenticate Request: Unable to generate a new token for user ", userID)
			return err
		}
		log.Info("JWT token refreshed for user ", userID)
		return nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			//Token is not valid
			return errors.New("Not a valid token!")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return errors.New("Token has expired or is not valid yet! Please login again to get a new token")
		} else {
			//Unknown error in the token
			return errors.New("Cannot verify this token! Please login again and get a new token!")
		}
	} else {
		//Unknown error in the token
		return errors.New("Cannot verify this token! Please login again and get a new token!")
	}
}

func AuthenticateServerToken(serv *Models.Server, token []byte) error {
	json_data, err := Utilities.DecryptAES(serv.Key, token)
	if err != nil {
		return err
	}
	type Data struct {
		Exp time.Time `json:"exp"`
		ID  string    `json:"id"`
	}
	var data Data
	err = json.Unmarshal(json_data, &data)
	if err != nil {
		log.Info("Server Token Authentication: Wrong format json, Error: ", err)
		return err
	}

	//Validate token expiry
	if data.Exp.Before(time.Now()) {
		log.Info("Server Token Authentication: Request token has expired. Server ID: ", serv.ServerID)
		return errors.New("Token is expired.")
	}

	//Validate token source
	if serv.ServerID != data.ID {
		log.Info("Server Token Authentication: Server IDs in token and request url does nto match")
		return errors.New("Token Source is un-authorized!")
	}

	//Everything looks fine.
	return nil
}
