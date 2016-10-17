package Server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/anshul35/ownit/Auth/JWT"
	"github.com/anshul35/ownit/Models"
	"github.com/anshul35/ownit/Utilities"

	log "github.com/golang/glog"
	"github.com/gorilla/mux"
)

const BasePath = "/server"

func DecodeJsonData(body io.ReadCloser) (map[string]interface{}, error) {
	var data interface{}
	decoder := json.NewDecoder(body)
	defer body.Close()
	err := decoder.Decode(&data)
	if err != nil {
		log.Error("URL handler: unable to decode json data. Error:", err)
		return nil, err
	}
	dict := data.(map[string]interface{})
	return dict, nil
}

var RegisterServerHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		serv := new(Models.Server)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(serv)
		if err != nil {
			log.Info("URL Handler: Wrong format json posted for registering server")
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Please POST json data in correct format"))
			return
		}
		defer serv.Save()
		token, err := GenerateServerToken()
		if err != nil {
			log.Error("URL Handler: Unable to generate server token while registering the server!")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Unable to register temporarily! Issue has been reported. Please try again after some time!"))
			return
		}
		type Response struct {
			Token string
		}
		resp := Response{Token: token}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error("URL Handler: Unable to Encode json response while registering server")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Unable to register temporarily! Issue has been reported. Please try again after some time!"))
		}
		log.Info("URL Handler: Successfully registered server: ", *serv)
		return
	})

var AddCommandHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		err := JWT.AuthenticateRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		comm := new(Models.Command)
		vars := mux.Vars(r)
		server, err := Models.GetServerByID(vars["serverID"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No server found for given server ID. Please re-register your server  with the cloud!"))
			return
		}
		dict, err := DecodeJsonData(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Unable to add command temporarily. Please try again after some time!"))
			return
		}
		defer r.Body.Close()
		for k, v := range dict {
			switch k {
			case "CommandString":
				comm.CommandString = v.(string)
				comm.CommandServer = server
				comm.CommandID = Utilities.GenerateUID()
			}
		}

		if comm != nil {
			defer comm.Save()
			log.Info("Successfully added command: ", comm.CommandID, ", ", comm.CommandString, " ; To server: ", server.ServerID)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("No command specified"))
			return
		}
	})

var ListCommandHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		var data = make([]Models.Command, 0)
		for _, v := range Models.CommandList {
			data = append(data, *v)
		}
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Unable to list commands temporarily. Please try again later!"))
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	})

var RunCommandHandler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		dict, err := DecodeJsonData(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Unable to run commands temporarily. Please try again later!"))
			return
		}
		defer r.Body.Close()

		requestID := ""
		commands := make([]*Models.Command, 0)
		errorComms := make([]string, 0)
		for k, v := range dict {
			switch k {
			case "CommandID":
				commList := v.([]interface{})
				for _, commID := range commList {
					comm, err := Models.GetCommandByID(commID.(string))

					//command ID does not exist
					if err != nil {
						errorComms = append(errorComms, commID.(string))
						continue
					}

					commands = append(commands, comm)
				}
				requestID = Utilities.GenerateUID()
				break
			}
		}

		//One or more command IDs do not exists
		if len(errorComms) != 0 {
			w.WriteHeader(http.StatusPartialContent)
			w.Write([]byte("No commands regisitered for command ID: " + strings.Join(errorComms, ",")))
			return
		}

		if len(commands) != 0 {
			_, err := Models.GetRequestByID(requestID)

			//Already existing request id
			if err == nil {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Request already exists"))
				return
			}

			runReq := Models.RunCommandRequest{
				RunCommands: commands,
				RequestID:   requestID,
			}

			err = runReq.Send()
			if err != nil {
				log.Error("URL Handler: Unable to send command run request. Error : ", err)
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("Unable to list commands temporarily. Please try again later!"))
				return
			}
			defer runReq.Save()

			type Response struct {
				RequestID string
			}
			resp := Response{RequestID: requestID}
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("Unable to list commands temporarily. Please try again later!"))
				return
			}
		} else {
			//No commands in the request
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Command list is empty. Please send atleast one command!"))
		}
		log.Info("URL Handler: Successfully responded to run command request!")
		return
	})
