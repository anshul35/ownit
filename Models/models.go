package Models

import "time"
 
type OwnList []Server

type Server struct {
	ServerID string
	IP string
	Name string
	Owner User
}

type Command struct {
	CommandID string
	CommandServer	Server
	CommandString	string
	Output	string
}

type User struct {
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

var ServerList []Server
var CommandList []Command

func All() []Gadget {
	return gadgetList
}