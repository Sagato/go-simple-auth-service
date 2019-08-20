package interfaces

import "authentication-service/model"

type Service interface {
	RegisterUser(u model.NewUser) error
}