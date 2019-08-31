package model

import "time"

type UserToken struct {
	RefreshToken string	`json:"refresh_token"`
	UserId int 			`json:"user_id"`
	Issued time.Time   	`json:"issued"`
	UserAgent string	`json:"user_agent"`
}
