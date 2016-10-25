package Models;

import (
	"errors"
//	"fmt"
	"net/http"
	"time"

	"github.com/anshul35/ownit/Settings/Constants"
)

type Server struct {
	ServerID string
	
	IP       string
	Name     string
	Owner    *User
	
	Token    string
	Key      []byte

	//True if Server.Start() has been called
	Active	bool

	//Server data output channels
	Processes chan []byte
	FileList chan []byte
}

var ServerList = make([]*Server, 0)

func NewServer() *Server {
	server := Server{Processes:make(chan []byte, Constants.ProcessesChannelSize),
		FileList:make(chan []byte, Constants.FileListChannelSize)}
	return &server
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

func (serv *Server) Save() bool {
	_, err := GetServerByID(serv.ServerID)
	if err != nil {
		ServerList = append(ServerList, serv)
		return true
	}
	return false
}

func (serv *Server) Start() {
	if serv.Active {
		return
	}
	
	client := &http.Client{Timeout:3*time.Second}
	port := Constants.ClientServerPort
	
	url := "http://" + serv.IP + ":" + port + "/api/v1/server/subscribe/processlist"
	client.Get(url)
	
	url = "http://" + serv.IP + ":" + port + "/api/v1/server/subscribe/filelist"
	client.Get(url)

	serv.Active = true
	return
}

func emptyChannel(c chan []byte) {
	n := len(c)
	for i := 0; i<n; i++ {
		<-c
	}
	return
}
func (serv *Server) Stop() {
	if !serv.Active {
		return
	}
	emptyChannel(serv.Processes)
	emptyChannel(serv.FileList)

	client := &http.Client{Timeout:3*time.Second}
	port := Constants.ClientServerPort
	
	url := "http://" + serv.IP + ":" + port + "/api/v1/server/unsubscribe/processlist"
	client.Get(url)
	
	url = "http://" + serv.IP + ":" + port + "/api/v1/server/unsubscribe/filelist"
	client.Get(url)

	serv.Active = false
	return
}
