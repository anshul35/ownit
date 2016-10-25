package Models

import (
	"errors"

	"github.com/anshul35/ownit/Utilities"
)

type Command struct {
	CommandID     string
	CommandServer *Server
	CommandString string
}

var CommandList = make([]*Command, 0)

func NewCommand() *Command {
	c := Command{CommandID:Utilities.GenerateUID()}
	return &c
}

type Object interface {
	Save() error
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

func (comm *Command) Run(requestId string, requestingUser *User) (*RunCommandRequest, error) {
	r := NewRunCommandRequest()
	r.RunCommands = []*Command{comm}
	r.RequestID = requestId
	r.RequestUser = requestingUser
	r.Server = comm.CommandServer
	r.Save()
	err := r.Send()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (comm *Command) Save() bool {
	_, err := GetCommandByID(comm.CommandID)
	if err != nil {
		CommandList = append(CommandList, comm)
		return true
	}
	return false
}