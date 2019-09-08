package api

import (
	"authentication-service/model"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (s *server) Token(w http.ResponseWriter, r *http.Request) {

	grantType, err := s.checkGrantType(r)

	if err != nil || grantType == "" {
		http.Error(w, "malformed request", http.StatusBadRequest)
		return
	}

	var user model.User

	err, code := s.resolveGrantTypeAndUser(&user, r, grantType);
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	blocked := s.checkIfUserBlocked(user)

	if blocked {
		http.Error(w, "your account is blocked, please contact support", http.StatusUnauthorized)
		return
	}

	tokenRes, err := s.generateTokenResponse(user)
	if err != nil {
		http.Error(w, "something went wrong. we are already investigating.", http.StatusInternalServerError)
		return
	}

	// add User Agent field to model to persist in db, but its omitted by json
	tokenRes.UserAgent = r.Header.Get("User-Agent")

	if err := s.db.Insert(&tokenRes); err != nil {
		http.Error(w, "something went wrong. we are already investigating", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tokenRes); err != nil {
		http.Error(w, "something went wrong. we are already investigating.", http.StatusInternalServerError)
	}

}

func (s *server) resolveGrantTypeAndUser(u *model.User, r *http.Request, grantType string) (error, int) {

	switch grantType {
	case "password": {
		var gt model.GrantTypePassword

		if err := s.decodeJson(r, &gt); err != nil {
			return errors.New("malformed request"), http.StatusBadRequest
		}

		if err := s.db.Model(u).Where("email = ?", gt.Username).Select(); err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return errors.New("unauthorized"), http.StatusUnauthorized
			}
		}

		err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(gt.Password))

		if subtle.ConstantTimeCompare([]byte(gt.Username), []byte(u.Email)) != 1 || err != nil {
			return errors.New("unauthorized"), http.StatusUnauthorized
		}
	}
	case "refresh_token": {
		var gt model.GrantTypeRefreshToken
		if err := s.decodeJson(r, &gt); err != nil {
			return errors.New("malformed request"), http.StatusBadRequest
		}

		var claims jwt.StandardClaims
		tkn, err := jwt.ParseWithClaims(gt.RefreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("jwtSecret")), nil
		})

		if err != nil {
			return errors.New(""), http.StatusUnauthorized
		}

		if !tkn.Valid {
			return errors.New(""), http.StatusUnauthorized
		}

		if err := s.db.Model(&u).Where("id = ?", claims.Subject).Select(); err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return errors.New(""), http.StatusUnauthorized
			}
		}
	}
	default:
		return errors.New( "no valid grant_type"), http.StatusBadRequest
	}

	return nil, 0
}

func (s *server) generateTokenResponse(user model.User) (model.GrantTypeResponse, error) {

	tokenExp := time.Now().Add(6 * time.Minute).Unix()
	refreshTokenExp := time.Now().Add(1 * time.Minute).Unix()

	uId := strconv.Itoa(user.Id)

	claims := &jwt.StandardClaims {
		Subject: uId,
		ExpiresAt: tokenExp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	rClaims := &jwt.StandardClaims {
		Subject: uId,
		ExpiresAt: refreshTokenExp,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)

	tokenString, err := token.SignedString([]byte(os.Getenv("jwtSecret")))
	if err != nil {
		return model.GrantTypeResponse{}, err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("jwtSecret")))
	if err != nil {
		return model.GrantTypeResponse{}, err
	}

	return model.GrantTypeResponse { UserId: user.Id, User: user, TokenType: "Bearer", ExpiresIn: tokenExp, AccessToken: tokenString, RefreshToken: refreshTokenString }, nil
}

func (s *server) checkGrantType(r *http.Request) (string, error) {
	jsonMap := make(map[string]string)
	if err := s.decodeJson(r, &jsonMap); err != nil {
		return "", err
	}
	return jsonMap["grant_type"], nil
}

func (s *server) checkIfUserBlocked(u model.User) bool {
	fmt.Println(u)
	return u.BlockedUntil.After(time.Now())
}