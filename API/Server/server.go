package Server;

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/anshul35/ownit/Models"
)

var BasePath = "/server"

var RegisterServerHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
		var serv Models.Server
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err := decoder.Decode(&serv)
		if err != nil {
			fmt.Println("cannot deocde the post request: ", err)
		}
		Models.ServerList = append(Models.ServerList, serv)
})

var AddCommandHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
	
})

var ListCommandHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
	
})

var RunCommandHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){

})