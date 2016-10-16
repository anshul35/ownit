package Models

import (
	"time"
	"fmt"
	"errors"
	"encoding/json"
	"bytes"
	"net/http"

	"github.com/anshul35/ownit/Settings/Constants"
)
 
type OwnList []Server

type Server struct {
	ServerID string
	IP string
	Name string
	Owner User
}

type Command struct {
	CommandID string
	CommandServer	*Server
	CommandString	string
}

type User struct {
	UserID int
	Servers OwnList
	Name     string
	Age      int
	FB       FacebookUser
	email    string
	birthday time.Time
}

type FacebookUser struct {
	FBID    int
	Handle  string
	Friends []FacebookUser
}

var ServerList = make([]*Server, 0)
var CommandList = make([]*Command, 0)

type Object interface {
	PopulateByID(string) (error)
	Save() (error)
}

type RunCommandRequest struct {
	RunCommands []*Command
	RequestID string
	RequestUser *User
	Output string
}

var ActiveRequests = make(map[string]*RunCommandRequest)
var OldRequests = make(map[string]*RunCommandRequest)

func GetRequestByID(id string) (*RunCommandRequest, error){
	if req, ok := ActiveRequests[id]; ok {
		return req, nil
	}else {
		return nil, errors.New("No request corresponding to given request id")
	}
}

func GetCommandByID(id string) (comm *Command, err error){
	for _,v := range CommandList {
		if v.CommandID == id {
			return v, nil
		}
	}
	err = errors.New("No server found")
	return nil, err
}

func GetServerByID(id string) (serv *Server, err error){
	for _,v := range ServerList {
		if v.ServerID == id {
			return v, nil
		}
	}
	err = errors.New("No server found")
	return nil, err
}

func (comm *Command) Save() bool {
	_, err := GetCommandByID(comm.CommandID)
	if (err != nil) {
		CommandList = append(CommandList, comm)
		return true
	}
	return false
}

func (serv *Server) Save() bool {
	_, err := GetServerByID(serv.ServerID)
	if (err != nil) {
		ServerList = append(ServerList, serv)
		return true
	}
	return false
}

func (req *RunCommandRequest) Save() bool {
	_, err := GetRequestByID(req.RequestID)
	if (err != nil) {
		ActiveRequests[req.RequestID] = req
		return true
	}
	return false
}


func (req *RunCommandRequest) Send() error {
	type Data struct {
		cloudToken string
		requestID string
		commands []string
	}
	comms := make([]string, 0)
	var server = new(Server)
	for _,v := range req.RunCommands {
		comms = append(comms, v.CommandString)
		server = v.CommandServer
	}

	//build json data for run request
	data := Data{requestID : req.RequestID,commands : comms}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)

	//build server url
	host := server.IP
	port  := Constants.ClientServerPort
	url := host + ":" + port + "/api/v1/server/commands/run"

	//Post data at url
	resp, _ := http.Post(url, "application/json; charset=utf-8", b)
	fmt.Println(resp)
	return nil
}