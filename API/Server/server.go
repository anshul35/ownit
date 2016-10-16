package Server;

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"

	"github.com/anshul35/ownit/Models"
	"github.com/anshul35/ownit/Utilities"

	"github.com/gorilla/mux"
)

var BasePath = "/server"

func DecodeJsonData(body io.ReadCloser) (map[string]interface{}, error) {
	var data interface{}
	decoder := json.NewDecoder(body)
	defer body.Close()
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	dict := data.(map[string]interface{})
	return dict, nil
}

var RegisterServerHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
		serv := new(Models.Server)
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(serv)
		if err != nil {
			fmt.Println("cannot deocde the post request: ", err)
			return
		}
		serv.Save()
		token, err := GenerateServerToken()
		if err != nil{
			fmt.Println("cannot generate token")
			return
		}
		type Response struct{
			Token string
		}
		resp := Response{Token:token}
		json.NewEncoder(w).Encode(resp)
		fmt.Println(Models.ServerList)
})

var AddCommandHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
		comm := new(Models.Command)
		vars := mux.Vars(r)
		server, err := Models.GetServerByID(vars["serverID"])
		if err != nil {
			fmt.Println(err)
			return
		}
		dict, err := DecodeJsonData(r.Body)
		defer r.Body.Close()
		for k, v := range dict {
			switch k {
				case "CommandString":
					comm.CommandString = v.(string)
					comm.CommandServer = server
					comm.CommandID = Utilities.GenerateUID()
			}
		}
		if err != nil {
			fmt.Println("cannot deocde the post request: ", err)
			return
		}
		if comm != nil {
			fmt.Println("comm is : ", comm)
			comm.Save()
		}
		fmt.Println(Models.CommandList)
})

var ListCommandHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
		var data = make([]Models.Command, 0)
		for _, v := range Models.CommandList {
			data = append(data, *v)
		}
		json.NewEncoder(w).Encode(data)
		return
})

var RunCommandHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
		dict, err := DecodeJsonData(r.Body)
		_ = err
		defer r.Body.Close()
		requestID := ""
		commands := make([]*Models.Command, 0)
		for k, v := range dict {
			switch k {
				case "CommandID":
					commList := v.([]interface{})
					for _, commID := range commList{
						comm, err := Models.GetCommandByID(commID.(string))
						if err != nil {
							fmt.Println("Error: ",err)
							continue
						}
						commands = append(commands, comm)
					}
					requestID = Utilities.GenerateUID()
					break
			}
		}
		if (requestID != "" && len(commands) != 0) {
			_, err := Models.GetRequestByID(requestID)
			if err == nil {
				//Already existing request id
				fmt.Println("Request id already exists!")
				return
			}
			runReq := Models.RunCommandRequest{
				RunCommands:commands,
				RequestID:requestID,
			}
			fmt.Println("Sending request now...")
			runReq.Send()
			fmt.Println("Sent! :)")
			runReq.Save()
			type Response struct{
				RequestID string
			}
			resp := Response{RequestID:requestID}
			json.NewEncoder(w).Encode(resp)
			return
		}
		return
})