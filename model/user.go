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
	Id int `json:"id" sql:,pk`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Active bool `json:"active sql:",notnull"`
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
