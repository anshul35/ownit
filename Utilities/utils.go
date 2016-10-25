package Utilities

import (
	"strconv"
	"time"
)

func init() {
	allowedChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for _, l := range allowedChars {
		ShortnerCharactors = append(ShortnerCharactors, string(l))
	}
}

func GenerateUID() string {
	curTime := int(time.Now().UnixNano())
	return strconv.Itoa(curTime)
}

var ShortnerCharactors = make([]string, 0)

func StringShortner(str string) (string, error) {
	//Input string should a number
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return "", err
	}

	base := int64(len(ShortnerCharactors))
	index := 0
	out := ""
	for {
		rem := id % base
		id = id / base
		out = ShortnerCharactors[rem] + out
		index = index + 1
		if id <= 0 {
			break
		}
	}
	return out, nil
}

func StringExpander(token string) string {
	base := int64(len(ShortnerCharactors))
	data := []rune(token)
	var it int64
	var res int64
	var val int64
	it = 1
	res = 0
	for i := len(data) - 1; i >= 0; i-- {
		if data[i]-'A' > 25 {
			val = int64(data[i] - 'a' + 26)
		} else {
			val = int64(data[i] - 'A')
		}
		res += it * val
		it *= base
	}
	return strconv.FormatInt(res, 10)
}
