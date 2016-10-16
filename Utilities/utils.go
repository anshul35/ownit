package Utilities;

import(
	"time"
	"strconv"
)


func GenerateUID() string{
	curTime := int(time.Now().UnixNano())
	return strconv.Itoa(curTime)
}

