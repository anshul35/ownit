package Models;

import(
	"errors"
)

type User struct {
	UserID   string
	Servers  []*Server
	Name     string
	FB       FacebookUser
	email    string
	JWTToken string
}

type FacebookUser struct {
	FBID int
	//	Handle  string
	//	Friends []FacebookUser
}

var UserList = make([]*User, 0)

func GetUserByID(id string) (serv *User, err error) {
	for _, v := range UserList {
		if v.UserID == id {
			return v, nil
		}
	}
	err = errors.New("No User found corresponding to user id:" + id)
	return nil, err
}

func (user *User) Save() bool {
	_, err := GetUserByID(user.UserID)
	if err != nil {
		UserList = append(UserList, user)
		return true
	}
	return false
}