package model

type AuthSession struct {
	AccessToken string 	`json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserId int 			`json:"user_id"`
	Scopes string 		`json:"scopes"`
}
