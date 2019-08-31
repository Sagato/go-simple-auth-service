package api

import (
	"authentication-service/db"
	"authentication-service/hashing/bcrypt"
	"authentication-service/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//func TestRegisterUser_Success(t *testing.T) {
//
//	dbConn := db.NewTestDB()
//	db.CreateTestSchemas(dbConn)
//
//
//	defer dbConn.Close()
//
//	s := NewServer()
//	s.router = mux.NewRouter()
//	s.hashing = bcrypt.New(4)
//	s.db = dbConn
//	emailConfig := email.SetupEmailCredentials(
//		os.Getenv("serverHost"),
//		os.Getenv("serverPort"),
//		os.Getenv("senderAddress"),
//		os.Getenv("username"),
//		os.Getenv("password"))
//
//	s.email = email.NewEmailSender(&emailConfig,"test")
//
//	testUser := model.NewUser{
//		Username: "Sagat",
//		Email: "danyal.iqbal@hotmail.com",
//		Password: "werwrwerwer",
//	}
//
//	reqBody, err := json.Marshal(testUser)
//
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//
//	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(s.RegisterUser)
//
//	handler.ServeHTTP(rr, req)
//
//	if status := rr.Code; status != http.StatusCreated {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusCreated)
//	}
//
//}

func TestRegisterUser_Exists(t *testing.T) {

	dbConn := db.NewTestDB()
	db.CreateTestSchemas(dbConn)

	defer dbConn.Close()

	s := NewServer()
	s.router = mux.NewRouter()
	s.hashing = bcrypt.New(4)
	s.db = dbConn

	testUser := model.User{
		Username: "Sagat",
		Email: "danyal.iqbal@hotmail.com",
		PasswordHash: "werwrweasadfsadfsadfsadfsadrwer",
		Active: false,
	}

	// Insert User into db and try to register same user
	if err := s.db.Insert(&testUser); err != nil {
		t.Errorf("Inserting test user went wrong %s", err.Error())
	}

	testUser2 := model.NewUser{
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

	s := NewServer()
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

	s := NewServer()
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

func TestUserSchema_Constraints(t *testing.T) {

	// Setting up test Environment
	dbConn := db.NewTestDB()
	db.CreateTestSchemas(dbConn)


	defer dbConn.Close()

	s := NewServer()
	s.router = mux.NewRouter()
	s.hashing = bcrypt.New(4)
	s.db = dbConn


	dbUser := &model.User{
		Username:     "Sagat",
		Email:        "",
		PasswordHash: "",
		Active:       false,
	}

	if err := s.db.Insert(dbUser); err != nil {
		t.Error(err.Error())
	}
}

func TestActivateAccount_MissingUrlParams(t *testing.T) {
	// Setting up test Environment
	dbConn := db.NewTestDB()
	db.CreateTestSchemas(dbConn)


	defer dbConn.Close()

	s := NewServer()
	s.router = mux.NewRouter()
	s.hashing = bcrypt.New(4)
	s.db = dbConn

	req, err := http.NewRequest("GET", "/activate", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ActivateAccount)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestActivateAccount_WrongToken(t *testing.T) {
	// Setting up test Environment
	dbConn := db.NewTestDB()
	if err := db.CreateTestSchemas(dbConn); err != nil {
		t.Fatal(err)
	}

	defer dbConn.Close()

	s := NewServer()
	s.router = mux.NewRouter()
	s.db = dbConn

	dbUser := model.User {
		Username:     "sagat",
		Email:        "sagat@hotm.com",
		PasswordHash: "asdiuhasdfiuh&76876asdjhg",
	}

	if err := s.db.Insert(&dbUser); err != nil {
		t.Error(err)
	}

	a := model.Activation{
		Token: "assdfhjkasddfkfjhasddfkfhj",
		Email: dbUser.Email,
	}

	if err := s.db.Insert(&a); err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", "/activate?token=ohjwfasdfhvasdfkjhbasdf", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ActivateAccount)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestActivation_EmptyEmailAddress(t *testing.T) {

	dbConn := db.NewTestDB()

	if err := db.CreateTestSchemas(dbConn); err != nil {
		t.Fatal("unable to create schema: ", err.Error())
	}

	s := &server{}
	s.db = dbConn
	s.router = mux.NewRouter()

	token := uuid.New().String()

	// Setup DB for test and insert needed test
	if err := s.db.Insert(&model.Activation{Token: token}); err != nil {
		t.Fatalf("db insertion failed: %s", err.Error())
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/activate?token=%s", token), nil)

	if err != nil {
		t.Errorf("request creation failed: %s", err.Error())
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.ActivateAccount)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}