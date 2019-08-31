package api

import (
	"authentication-service/model"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
	"time"
)

func (s *server) Token(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("User Agent: %s", r.Header.Get("User-Agent"))

	grantType, err := s.checkGrantType(r)

	if err != nil || grantType == "" {
		http.Error(w, "malformed request", http.StatusBadRequest)
		return
	}

	var user model.User

	switch grantType {

	case "password": {
		var gt model.GrantTypePassword

		if err := s.decodeJson(r, &gt); err != nil {
			http.Error(w, "malformed request", http.StatusBadRequest)
			return
		}

		if err := s.db.Model(&user).Where("email = ?", gt.Username).Select(); err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(gt.Password))

		if subtle.ConstantTimeCompare([]byte(gt.Username), []byte(user.Email)) != 1 || err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	case "refresh_token": {
		var gt model.GrantTypeRefreshToken
		if err := s.decodeJson(r, &gt); err != nil {
			http.Error(w, "malformed request", http.StatusBadRequest)
			return
		}

		var claims jwt.StandardClaims

		tkn, err := jwt.ParseWithClaims(gt.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
			return os.Getenv("jwtSecret"), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := s.db.Model(&user).Where("id = ?", claims.Subject).Select(); err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
	}

	default:
		http.Error(w, "no valid grant_type", http.StatusBadRequest)
		return
	}

	blocked := s.checkIfUserBlocked(user)

	if blocked {
		http.Error(w, "your account is blocked, please contact support", http.StatusUnauthorized)
		return
	}

	tokenRes, err := s.generateTokenResponse(user)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "something went wrong. we are already investigating.", http.StatusInternalServerError)
		return
	}

	userToken := model.UserToken{
		RefreshToken: tokenRes.RefreshToken,
		UserId:       user.Id,
		Issued:       time.Now().UTC(),
		UserAgent: 	  r.Header.Get("User-Agent"),
	}

	if err := s.db.Insert(&userToken); err != nil {
		http.Error(w, "something went wrong. we are already investigating", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tokenRes); err != nil {
		http.Error(w, "something went wrong. we are already investigating.", http.StatusInternalServerError)
	}

}

func (s *server) generateTokenResponse(user model.User) (model.GrantTypeResponse, error) {

	fmt.Printf("exp time: %d", time.Now().Add(6 * time.Minute).Unix())

	tokenExp := time.Now().Add(6 * time.Minute).Unix()
	refreshTokenExp := time.Now().Add(10 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
		"sub": user.Id,
		"exp": tokenExp,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": refreshTokenExp,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("jwtSecret")))
	if err != nil {
		fmt.Println(err.Error())
		return model.GrantTypeResponse{}, err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte("test"))
	if err != nil {
		fmt.Println(err.Error())
		return model.GrantTypeResponse{}, err
	}

	return model.GrantTypeResponse { TokenType: "Bearer", ExpiresIn: tokenExp, AccessToken: tokenString, RefreshToken: refreshTokenString }, nil
}

func (s *server) checkGrantType(r *http.Request) (string, error) {
	jsonMap := make(map[string]string)
	if err := s.decodeJson(r, &jsonMap); err != nil {
		return "", err
	}
	fmt.Printf("Json Map: %s", jsonMap)
	return jsonMap["grant_type"], nil
}

func (s *server) checkIfUserBlocked(u model.User) bool {
	return u.BlockedUntil.After(time.Now())
}