package api

import (
	"authentication-service/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *server) RegisterUser(w http.ResponseWriter, r *http.Request) {

	// Incoming newUser
	var newUser model.NewUser

	// Decode the body on the request
	if err := s.decodeJson(r, &newUser); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if newUser already exists
	var dbUser model.User

	exists, err := s.db.Model(&dbUser).Where("email = ?", newUser.Email).WhereOr("username = ?", newUser.Username).Exists()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		err := errors.New("user already exists")
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
		fmt.Println(err.Error())
		http.Error(w, "Something went wrong. We are sorry and already investigating the issue.", http.StatusInternalServerError)
		return
	}

	uuid := uuid.New()

	activation := model.Activation{
		Token: uuid.String(),
		Email: newUser.Email,
	}

	if err := s.db.Insert(&activation); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Init Anonymous struct with template data
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
		dbUser.Username,
		"Hello",
		"https://www.example.com",
		dbUser.Email,
		"Sagateyson",
		"2019",
		"sagat corp.",
		"https://www.example.com/activate",
		"testdomain.com",
	}

	wd, err := os.Getwd()

	path, err := filepath.Abs(wd)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Path: \n %s", path)

	// Get working Directory for correct html template file paths
	if err != nil {
		http.Error(w, "Something went wrong. We are already investigating", http.StatusInternalServerError)
		return
	}

	if err := s.email.ParseTemplate(filepath.Join(wd,  "./email/templates/registration_mail.html"), data); err != nil {
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

func (s *server) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	token, ok := r.URL.Query()["token"]

	if !ok || len(token) < 1 {
		log.Println("url params are missing")
		http.Error(w, "url params are missing", http.StatusBadRequest)
		return
	}

	fmt.Println(token[0])

	activation := model.Activation{}

	if err := s.db.Model(&activation).Where("token = ? ", token[0]).OnConflict("token", "DO NOTHING").Select(); err != nil {

		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "no such account", http.StatusNotFound)
			return
		}

		http.Error(w, "something went wrong. we are already investigating", http.StatusInternalServerError)
		return
	}

	if activation.IsEmpty() {
		http.Error(w, "no such account", http.StatusNotFound)
		return
	}

	var user model.User

	if err := s.db.Model(&user).Where("email = ?", activation.Email).Select(); err != nil {

		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "no such account", http.StatusNotFound)
			return
		}

		http.Error(w, "something went wrong. we are already investigating", http.StatusInternalServerError)
		return
	}

	if user.IsEmpty() {
		http.Error(w, "no such user", http.StatusNotFound)
		return
	}

	user.Active = true

	_, err := s.db.Model(&user).Set("active = ?active").Where("id = ?id").Update()

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "something went wrong. we are already investigating", http.StatusInternalServerError)
		return
	}

	if err := s.db.Delete(&activation); err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
