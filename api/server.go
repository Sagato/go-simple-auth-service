package api

import (
	"authentication-service/db"
	"authentication-service/email"
	"authentication-service/hashing"
	"authentication-service/hashing/bcrypt"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

func (s *server) decodeJson(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Req Body: %v", string(body))

	if err := json.Unmarshal(body, &v); err != nil {
		return err
	}

	// reset the response body to initial state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return nil
}

func (s *server) encodeJson(v interface{}) ([]byte, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}
	return data, nil
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

	port, err := strconv.Atoi(os.Getenv("serverPort"))
	if err != nil {
		panic(err.Error())
	}

	emailConfig := email.SetupEmailCredentials(
		os.Getenv("serverHost"),
		os.Getenv("senderAddress"),
		os.Getenv("username"),
		os.Getenv("password"),
		port,
	)

	es := email.NewEmailSender(&emailConfig,	"test", func(to []string) {})

	s := &server {conn, r, bcrypt.New(4), es }
	s.routes()

	if err := http.ListenAndServe(":1234", s.router); err != nil {
		log.Fatal(err.Error())
	}
}
