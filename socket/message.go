package socket;

import(
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anshul35/ownit/Models"
)

const(
	TypeProcessList = "PL"
	TypeFileList = "FL"
	TypeCommandOutput = "CO"
)

type Message struct {
	isInbound bool `json:"-"`

	client *Client `json:"-"`

	ID string `json:"id"`

	Type string `json:"type"`

	Body []byte `json:"body"`

	ServerID string `json:"server_id"`
}


//Converts msg struct to json string
func (m *Message) encode() ([]byte, error) {
	s, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *Message) Execute(attachOutput bool) (*Models.RunCommandRequest, error){
	_ = attachOutput
	if (m.Type != TypeCommandOutput) || !(m.isInbound){
		fmt.Println("Error in execution")
		return nil, errors.New("Cannot execute given messsage type!")
	}
	s, err := Models.GetServerByID(m.ServerID)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	c := Models.NewCommand()
	c.CommandServer = s
	c.CommandString = string(m.Body)
	c.Save()
	r, err := c.Run(m.ID, m.client.User)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	m.client.registerRequest(r)
	return r, err
}