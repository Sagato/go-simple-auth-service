package api

import (
	"authentication-service/model"
	"testing"
)

func TestJwt_Created(t *testing.T) {
	s := &server{}

	jwtRes, err := s.generateTokenResponse(model.User{Id: 1234})
	if err != nil {
		t.Errorf("generating jwt went wrong: %v", jwtRes)
	}

	if jwtRes.RefreshToken == "" || jwtRes.AccessToken == "" {
		t.Errorf("generating jwt went wrong: %v", jwtRes)
	}
}
