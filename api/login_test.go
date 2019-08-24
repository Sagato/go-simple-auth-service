package api

import (
	"authentication-service/db"
	"authentication-service/model"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin_NoCredentials(t *testing.T) {

	s := &server{}
	initDBAndCreateTestSchema(s,t)
	s.router = mux.NewRouter()

	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Errorf("generating request went wrong: %s", err.Error())
	}

	req.SetBasicAuth("", "")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("returned wrong status code: got %d, want %d", status, http.StatusUnauthorized)
	}
}

func TestLogin_WrongCredentials(t *testing.T) {

	s := &server{}
	initDBAndCreateTestSchema(s,t)
	s.router = mux.NewRouter()

	hash, err := bcrypt.GenerateFromPassword([]byte("123sdfiuhdfgoiuwhfdgoi&(/&456"), 4)
	if err != nil {
		t.Error(err)
	}

	user := &model.User {
		Id: 0,
		Username: "sagat",
		Email: "danyal.iqbal@danyal.com",
		PasswordHash: string(hash),
		Active: true,
	}

	if err := s.db.Insert(user); err != nil {
		t.Errorf("insert into db went wring: %s", err.Error())
	}

	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Errorf("generating request went wrong: %s", err.Error())
	}

	req.SetBasicAuth(user.Email, "1asdfafasdfasf23424%&/")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("returned wrong status code: got %d, want %d", status, http.StatusUnauthorized)
	}
}

func initDBAndCreateTestSchema(s *server, t *testing.T) {
	dbConn := db.NewTestDB()
	if err := db.CreateTestSchemas(dbConn); err != nil {
		t.Errorf("init db went wrong: %s", err.Error())
	}
	s.db = dbConn
}
