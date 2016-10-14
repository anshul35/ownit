package Models

import (
	"time"
	"errors"
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
	Output	string
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

type Gadget struct {
	DisplayName string
	Make        string
	Model       string
	ReleaseDate time.Time
}

var gadgetList = []Gadget{
		Gadget{"Apple Watch", "Apple", "Watch 2.0", time.Now()},
		Gadget{"Moto 360", "Moto", "1.0", time.Now()}}

var ServerList []Server = make([]Server, 1, 2)
var CommandList []Command = make([]Command, 2, 5)

func GetServerByID(id string) (err error, res *Server) {
	for _,v := range ServerList {
		if v.ServerID == id {
			res = &v
			return
		}
	}
	err = errors.New("No server found")
	return 
}

func All() []Gadget {
	return gadgetList
}