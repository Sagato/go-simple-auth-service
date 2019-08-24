package api

import (
	"authentication-service/model"
	"crypto/subtle"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func (s *server) Login(w http.ResponseWriter, r *http.Request) {

	var user model.User

	username, pw, ok := r.BasicAuth()

	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="login"`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := s.db.Model(&user).Where("email = ?", username).Select(); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(pw))

	if !ok || subtle.ConstantTimeCompare([]byte(username), []byte(user.Email)) != 1 || err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

}
