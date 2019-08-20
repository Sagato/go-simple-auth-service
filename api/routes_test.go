package api

import (
	"authentication-service/db"
	"authentication-service/email"
	"authentication-service/hashing/bcrypt"
	"authentication-service/model"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Test Register User endpoint in success case
func TestRegisterUser_Success(t *testing.T) {

	dbConn := db.NewTestDB()
	db.CreateTestSchemas(dbConn)


	defer dbConn.Close()

	s := newServer()
	s.router = mux.NewRouter()
	s.hashing = bcrypt.New(4)
	s.db = dbConn
	emailConfig := email.SetupEmailCredentials(
		os.Getenv("serverHost"),
		os.Getenv("serverPort"),
		os.Getenv("senderAddress"),
		os.Getenv("username"),
		os.Getenv("password"))

	s.email = email.NewEmailSender(&emailConfig,"test")

	testUser := model.NewUser{
		Username: "Sagat",
		Email: "danyal.iqbal@hotmail.com",
		Password: "werwrwerwer",
	}

	reqBody, err := json.Marshal(testUser)

	if err != nil {
		t.Errorf(err.Error())
	}

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.RegisterUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

}

func TestRegisterUser_Exists(t *testing.T) {

	dbConn := db.NewTestDB()
	db.CreateTestSchemas(dbConn)

	defer dbConn.Close()

	s := newServer()
	s.router = mux.NewRouter()
	s.hashing = bcrypt.New(4)
	s.db = dbConn

	testUser := model.User {
		Username: "Sagat",
		Email: "danyal.iqbal@hotmail.com",
		PasswordHash: "werwrweasadfsadfsadfsadfsadrwer",
		Active: false,
	}

	// Insert User into db and try to register same user
	if err := s.db.Insert(&testUser); err != nil {
		t.Errorf("Inserting test user went wrong %s", err.Error())
	}

	testUser2 := model.NewUser {
		Username: "Sagat",
		Email: "danyal.iqbal@hotmail.com",
		Password: "werwrweasadfsadfsadfsadfsadrwer",
	}

	reqBody, err := json.Marshal(testUser2)

	if err != nil {
		t.Errorf(err.Error())
	}

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.RegisterUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK && !strings.Contains(rr.Body.String(), "user already exists") {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestRegisterUser_WrongFieldTypes(t *testing.T) {

	// Setting up test Environment
	dbConn := db.NewTestDB()
	db.CreateTestSchemas(dbConn)


	defer dbConn.Close()

	s := newServer()
	s.router = mux.NewRouter()
	s.hashing = bcrypt.New(4)
	s.db = dbConn


	testUser :=
		`{ 
		Username: false,
		Email: "danyal.iqbal@hotmail.com",
		Password: "werwrweasadfsadfsadfsadfsadrwer",
	}`

	reqBody, err := json.Marshal(testUser)

	if err != nil {
		t.Errorf(err.Error())
	}

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.RegisterUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

func TestRegisterUser_MalformedJson(t *testing.T) {

	// Setting up test Environment
	dbConn := db.NewTestDB()
	db.CreateTestSchemas(dbConn)


	defer dbConn.Close()

	s := newServer()
	s.router = mux.NewRouter()
	s.hashing = bcrypt.New(4)
	s.db = dbConn


	testUser :=
		`{ 
		Username: false,
		Email: "danyal.iqbal@hotmail.com",
		Password: "werwrweasadfsadfsadfsadfsadrwer",
	`

	reqBody, err := json.Marshal(testUser)

	if err != nil {
		t.Errorf(err.Error())
	}

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.RegisterUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}

