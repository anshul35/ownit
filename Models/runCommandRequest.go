package Models;

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/golang/glog"

	"github.com/anshul35/ownit/Settings/Constants"
)

type RunCommandRequest struct {
	RunCommands []*Command
	RequestID   string
	RequestUser *User
	Output      []byte
	OutChan 	chan []byte
	Server 		*Server
}

var ActiveRequests = make(map[string]*RunCommandRequest)
var OldRequests = make(map[string]*RunCommandRequest)

//**Note: Returns without saving
func NewRunCommandRequest() *RunCommandRequest {
	r := RunCommandRequest{OutChan:make(chan []byte, Constants.RequestChannelSize)}
	return &r
}

func GetRequestByID(id string) (*RunCommandRequest, error) {
	if req, ok := ActiveRequests[id]; ok {
		return req, nil
	} else {
		return nil, errors.New("No Run Request found corresponding to request id: " + id)
	}
}
func (req *RunCommandRequest) Save() bool {
	_, err := GetRequestByID(req.RequestID)
	if err != nil {
		ActiveRequests[req.RequestID] = req
		return true
	}
	return false
}

func (req *RunCommandRequest) Send() error {
	type Data struct {
		RequestID  string
		Commands   []string
	}

	comms := make([]string, 0)
	var server = new(Server)
	for _, v := range req.RunCommands {
		comms = append(comms, v.CommandString)
		server = v.CommandServer
	}

	//build json data for run request
	data := Data{RequestID: req.RequestID, Commands: comms}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		log.Error("Models: Error in encoding json. Error: ", err)
		return err
	}

	//build server url
	host := server.IP
	port := Constants.ClientServerPort
	url := "http://" + host + ":" + port + "/api/v1/server/commands/run"

	//Post data at url
	r, err := http.NewRequest("POST", url,  b)
	if err != nil {
		log.Error("Models: Error in creating Post request to run commands to client server. Error: ", err)
		return err
	}
	r.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout:3*time.Second}
	_, err = client.Do(r)
	if err != nil {
		log.Error("Models: Error in sending Post request to run commands to client server. Error: ", err)
		return err
	}
	log.Info("Models: Sucessfully sent POST request to client server: ", host, " ,to run commands: ", data.Commands)
	return nil
}
