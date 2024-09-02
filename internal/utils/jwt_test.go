package utils

import "testing"

func TestJWT(t *testing.T) {
	jwtService := NewJWT("secret", 10)
	token, err := jwtService.GenerateTokenWithGroup(1, "admin")
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("token is empty")
	}
	t.Log(token)

	claims, err := jwtService.ParseToken(token)
	if err != nil {
		t.Error(err)
	}
	if claims == nil {
		t.Error("claims is nil")
	}
	t.Log(claims)
}
