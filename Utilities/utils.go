package Utilities

import (
	"strconv"
	"time"
)

func GenerateUID() string {
	curTime := int(time.Now().UnixNano())
	return strconv.Itoa(curTime)
}
