package interfaces

import "authentication-service/model"

type Repository interface {
	RegisterUser(u model.NewUser) error
}
