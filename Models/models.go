package Models

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/golang/glog"

	"github.com/anshul35/ownit/Settings/Constants"
)

type OwnList []Server

type Server struct {
	ServerID string
	IP       string
	Name     string
	Owner    User
	Token    string
	Key      []byte
}

type Command struct {
	CommandID     string
	CommandServer *Server
	CommandString string
}

type User struct {
	UserID   string
	Servers  OwnList
	Name     string
	FB       FacebookUser
	email    string
	JWTToken string
}

type FacebookUser struct {
	FBID int
	//	Handle  string
	//	Friends []FacebookUser
}

var ServerList = make([]*Server, 0)
var CommandList = make([]*Command, 0)
var UserList = make([]*User, 0)

type Object interface {
	Save() error
}

type RunCommandRequest struct {
	RunCommands []*Command
	RequestID   string
	RequestUser *User
	Output      string
}

var ActiveRequests = make(map[string]*RunCommandRequest)
var OldRequests = make(map[string]*RunCommandRequest)

func GetRequestByID(id string) (*RunCommandRequest, error) {
	if req, ok := ActiveRequests[id]; ok {
		return req, nil
	} else {
		return nil, errors.New("No Run Request found corresponding to request id: " + id)
	}
}

func GetCommandByID(id string) (comm *Command, err error) {
	for _, v := range CommandList {
		if v.CommandID == id {
			return v, nil
		}
	}
	err = errors.New("No Command corresponding to command id: " + id)
	return nil, err
}

func GetServerByID(id string) (serv *Server, err error) {
	for _, v := range ServerList {
		if v.ServerID == id {
			return v, nil
		}
	}
	err = errors.New("No Server found corresponding to server id:" + id)
	return nil, err
}

func GetUserByID(id string) (serv *User, err error) {
	for _, v := range UserList {
		if v.UserID == id {
			return v, nil
		}
	}
	err = errors.New("No Server found corresponding to server id:" + id)
	return nil, err
}

func (comm *Command) Save() bool {
	_, err := GetCommandByID(comm.CommandID)
	if err != nil {
		CommandList = append(CommandList, comm)
		return true
	}
	return false
}

func (serv *Server) Save() bool {
	_, err := GetServerByID(serv.ServerID)
	if err != nil {
		ServerList = append(ServerList, serv)
		return true
	}
	return false
}

func (user *User) Save() bool {
	_, err := GetServerByID(user.UserID)
	if err != nil {
		UserList = append(UserList, user)
		return true
	}
	return false
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
		CloudToken string
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
	_, err = http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		log.Error("Models: Error in sending Post request to run commands to client server. Error: ", err)
		return err
	}
	log.Info("Models: Sucessfully sent POST request to client server: ", host, " ,to run commands: ", data.Commands)
	return nil
}
