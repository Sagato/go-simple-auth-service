package api

import (
	"authentication-service/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
)

func (s *server) RegisterUser(w http.ResponseWriter, r *http.Request) {

	// Incoming newUser
	var newUser model.NewUser

	// Decode the body on the request
	if err := s.decode(w, r, &newUser); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if newUser already exists
	var dbUser model.User

	exists, err := s.db.Model(&dbUser).Where("email = ?", newUser.Email).Exists()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		err := errors.New("User already exists")
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	// If not exists generate a Hash from the password
	pwHash, err := s.hashing.HashPassword(newUser.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	insertableUser := model.User{
		Id:           0,
		Username:     newUser.Username,
		Email:        newUser.Email,
		PasswordHash: pwHash,
		Active:       false,
	}

	if err := s.db.Insert(&insertableUser); err != nil {
		http.Error(w, "Something went wrong. We are sorry and already investigating the issue.", http.StatusInternalServerError)
		return
	}

	uuid := uuid.New()

	activation := model.Activation {
		Token: uuid.String(),
		Email: newUser.Email,
	}

	if err := s.db.Insert(&activation); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wd, err := os.Getwd()

	if err != nil {
		http.Error(w, "Something went wrong. We are already investigating", http.StatusInternalServerError)
		return
	}

	data := struct {
		Name string
		Greetings string
		WebsiteUrl string
		Email string
		WebsiteName string
		Year string
		CompanyName string
		ActivationLink string
		From string
	}{
		"sagat",
		"Hello",
		"https://www.example.com",
		"sagat@web.de",
		"Sagateyson",
		"2019",
		"sagat corp.",
		"https://www.example.com/activate",
		"testdomain.com",
	}

	if err := s.email.ParseTemplate(wd + "/email/templates/registration_mail.html", data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong. We are sorry and already investigating the issue.", http.StatusInternalServerError)
		return
	}

	if err := s.email.Send([]string{ newUser.Email }); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong. We are sorry and already investigating the issue.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
