package model

import (
	"encoding/json"
	"fmt"
)

type NewUser struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id int `json:"id" sql:",pk"`
	Username string `json:"username" sql:",notnull,unique"`
	Email string `json:"email" sql:",notnull,unique"`
	PasswordHash string `json:"password_hash" sql:",notnull"`
	Active bool `json:"active" sql:",notnull"`
}

type UserDTO struct {
	Username string
	Email string
}

func (u User) String() string {
	return fmt.Sprintf("User<%s %s>", u.Username, u.Email)
}

func (u User) MarshalJSON() ([]byte, error) {
	type user User
	x := user(u)
	x.PasswordHash = ""
	return json.Marshal(x)
}

func (u User) IsEmpty() bool {
	return u.Id == 0 && u.Username == "" && u.Email == "" && u.PasswordHash == "" && u.Active == false
}
