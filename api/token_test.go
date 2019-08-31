package api

import (
	"authentication-service/db"
	"authentication-service/model"
	"bytes"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogin_NoCredentials(t *testing.T) {

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

	gt := &model.GrantTypePassword{
		GrantType: "password",
		Username: "",
		Password: "",
	}

	data, err := s.encodeJson(gt)
	if err != nil {
		t.Errorf("marshalling json went wrong: %s", err.Error())
	}

	req, err := http.NewRequest("POST", "/accessToken", bytes.NewReader(data))
	if err != nil {
		t.Errorf("generating request went wrong: %s", err.Error())
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.Token)

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

	gt := &model.GrantTypePassword{
		GrantType: "password",
		Username: "danyal.iqbal@danyal.com",
		Password: "1asdfafasdfasf23424%&/",
	}

	data, err := s.encodeJson(gt)
	if err != nil {
		t.Errorf("marshalling json went wrong: %s", err.Error())
	}

	req, err := http.NewRequest("POST", "/accessToken",  bytes.NewReader(data))
	if err != nil {
		t.Errorf("generating request went wrong: %s", err.Error())
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.Token)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("returned wrong status code: got %d, want %d", status, http.StatusUnauthorized)
	}
}

func Test_CheckGrantTypePassword(t *testing.T) {
	gtPw := model.GrantTypePassword{
		"password",
		"sagat",
		"ojhabsdfokjabsdf",
	}

	s := &server{}

	reqBody, err := s.encodeJson(&gtPw)
	if err != nil {
		t.Errorf("encoding went wrong: %s", err.Error())
	}

	req, err := http.NewRequest("POST", "/token", bytes.NewReader(reqBody))

	grantType, err := s.checkGrantType(req)
	if err != nil {
		t.Errorf("checking grant type went wrong: %s", err.Error())
	}

	if grantType != "password" {
		t.Errorf("wrong grant type detected: %s", grantType)
	}

}

func Test_GrantTypeRefreshToken(t *testing.T) {
	gtRt := model.GrantTypeRefreshToken{
		"refresh_token",
		"sadfoijasdbfokpjnsadfkojpnasdfokjn",
	}

	s := &server{}

	reqBody, err := s.encodeJson(&gtRt)
	if err != nil {
		t.Errorf("encoding went wrong: %s", err.Error())
	}

	req, err := http.NewRequest("POST", "/token", bytes.NewReader(reqBody))

	grantType, err := s.checkGrantType(req)
	if err != nil {
		t.Errorf("checking grant type went wrong: %s", err.Error())
	}

	if grantType != "refresh_token" {
		t.Errorf("wrong grant type detected: %s", grantType)
	}
}

func Test_WrongGrantTypeHttp(t *testing.T) {
	gtPw := model.GrantTypePassword{
		"sdafka",
		"",
		"ojhabsdfokjabsdf",
	}

	s := &server{}

	reqBody, err := s.encodeJson(&gtPw)
	if err != nil {
		t.Errorf("encoding went wrong: %s", err.Error())
	}

	req, err := http.NewRequest("POST", "/token", bytes.NewReader(reqBody))
	if err != nil {
		t.Errorf("generating request went wrong: %s", err.Error())
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.Token)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("returned wrong status code: got %d, want %d", status, http.StatusBadRequest)
	}
}

func Test_NotBlockedUser(t *testing.T) {

	hash, err := bcrypt.GenerateFromPassword([]byte("123sdfiuhdfgoiuwhfdgo6"), 4)
	if err != nil {
		t.Error(err)
	}

	u := &model.User {
		Id: 0,
		Username: "sagat",
		Email: "danyal.iqbal@hotmail.com",
		PasswordHash: string(hash),
		Active: true,
		BlockedUntil: time.Now().AddDate(0, 0, -1),
	}

	s := &server{}
	initDBAndCreateTestSchema(s,t)
	s.router = mux.NewRouter()

	if err := s.db.Insert(u); err != nil {
		t.Errorf("insert into db went wring: %s", err.Error())
	}

	gt := &model.GrantTypePassword{
		GrantType: "password",
		Username: "danyal.iqbal@hotmail.com",
		Password: "123sdfiuhdfgoiuwhfdgo6",
	}

	data, err := s.encodeJson(gt)
	if err != nil {
		t.Errorf("marshalling json went wrong: %s", err.Error())
	}

	req, err := http.NewRequest("POST", "/accessToken", bytes.NewReader(data))
	if err != nil {
		t.Errorf("generating request went wrong: %s", err.Error())
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.Token)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("returned wrong status code: got %d, want %d", status, http.StatusOK)
	}

}

func Test_BlockedUser(t *testing.T) {

	hash, err := bcrypt.GenerateFromPassword([]byte("123sdfiuhdfgoiuwhfdgo6"), 4)
	if err != nil {
		t.Error(err)
	}

	u := &model.User {
		Id: 0,
		Username: "sagat",
		Email: "danyal.iqbal@hotmail.com",
		PasswordHash: string(hash),
		Active: true,
		BlockedUntil: time.Now().AddDate(0, 0, +1),
	}

	s := &server{}
	initDBAndCreateTestSchema(s,t)
	s.router = mux.NewRouter()

	if err := s.db.Insert(u); err != nil {
		t.Errorf("insert into db went wring: %s", err.Error())
	}

	gt := &model.GrantTypePassword{
		GrantType: "password",
		Username: "danyal.iqbal@hotmail.com",
		Password: "123sdfiuhdfgoiuwhfdgo6",
	}

	data, err := s.encodeJson(gt)
	if err != nil {
		t.Errorf("marshalling json went wrong: %s", err.Error())
	}

	req, err := http.NewRequest("POST", "/accessToken", bytes.NewReader(data))
	if err != nil {
		t.Errorf("generating request went wrong: %s", err.Error())
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.Token)

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
