package Server;

import (
	"fmt"
	"net/http"
	"reflect"
	"encoding/json"
	"github.com/anshul35/ownit/Models"

	"github.com/gorilla/mux"
)

var BasePath = "/server"

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
		Models.ServerList = append(Models.ServerList, *serv)
})

var AddCommandHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
		comm := new(Models.Command)
		var data interface{}
		vars := mux.Vars(r)
		err, server := Models.GetServerByID(vars["serverID"])
		if err != nil {
			fmt.Println(err)
			return
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err = decoder.Decode(&data)
		dict := data.(map[string]interface{})
		for k, v := range dict {
			switch k {
				case "CommandString":
					comm.CommandString = v.(string)
					comm.CommandServer = server
			}
		}
		if err != nil {
			fmt.Println("cannot deocde the post request: ", err)
			return
		}
		if comm != nil {
			Models.CommandList = append(Models.CommandList, *comm)
		}	
})

var ListCommandHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
	
})

var RunCommandHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){

})