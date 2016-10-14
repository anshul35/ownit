package Settings

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Facebook map[string]string
}

var Config Configuration

func init() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("fwqf", Config.Facebook) // output: [UserA, UserB]

}
