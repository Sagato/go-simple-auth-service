package model

import "time"

type GrantTypePassword struct {
	GrantType string 	`json:"grant_type"`
	Username string 	`json:"username"`
	Password string 	`json:"password"`
}

type GrantTypeRefreshToken struct {
	GrantType string 	`json:"grant_type"`
	RefreshToken string 	`json:"refresh_token"`
}

type GrantTypeResponse struct {
	Id int `json:"-" sql=",pk"`
	UserId int
	User User `json:"-"`
	TokenType string `json:"token_type"`
	ExpiresIn int64 `json:"expires_in"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserAgent string `json:"-"`
	CreatedAt time.Time `json:"-" sql:"default:now()"`
	UpdatedAt time.Time `json:"-"`
}
