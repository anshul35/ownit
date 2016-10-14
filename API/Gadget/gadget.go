package Gadget;

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/anshul35/ownit/Models"
)

var BasePath = "/gadget"

var GadgetListHandler = http.HandlerFunc(
	func (w http.ResponseWriter, r *http.Request){
		var gadgets, err = json.Marshal(Models.All())
		if err != nil {
			fmt.Println("Gadget list can not be prepared now bcs ", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(gadgets))
})