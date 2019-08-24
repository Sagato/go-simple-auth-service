package api

import (
	"authentication-service/db"
	"authentication-service/email"
	"authentication-service/hashing"
	"authentication-service/hashing/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type server struct {
	db      *pg.DB
	router  *mux.Router
	hashing hashing.Hashing
	email   email.Sender
}

func NewServer() *server {
	return &server{}
}

func (s *server) ServeHttp(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	r := mux.NewRouter()

	conn := db.NewDB()
	if err := db.CreateSchema(conn); err != nil {
		fmt.Println(err.Error())
	}

	emailConfig := email.SetupEmailCredentials(
		os.Getenv("serverHost"),
		os.Getenv("serverPort"),
		os.Getenv("senderAddress"),
		os.Getenv("username"),
		os.Getenv("password"))

	es := email.NewEmailSender(&emailConfig,	"test")

	s := &server {conn, r, bcrypt.New(4), es }
	s.routes()

	if err := http.ListenAndServe(":1234", s.router); err != nil {
		log.Fatal(err.Error())
	}
}
