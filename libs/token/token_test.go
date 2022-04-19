package token

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestValidateToken(t *testing.T) {
	svc := NewService("secret-key")
	claim := AuthServiceJWT{
		StandardClaims: jwt.StandardClaims{
			Id:        "uniqueID",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token, err := svc.Generate(claim)
	if err != nil {
		t.Error("validate token", err)
		return
	}

	respToken, err := svc.Validate(token)
	if err != nil {
		t.Error("validate token", err)
		return
	}

	claims := respToken.Claims.(jwt.MapClaims)
	byteClaims, err := json.Marshal(claims)
	if err != nil {
		t.Error(err)
		return
	}

	myClaim := AuthServiceJWT{}
	err = json.Unmarshal(byteClaims, &myClaim)
	if err != nil {
		t.Error(err)
		return
	}

	if myClaim.Id != claim.Id {
		t.Error("Id not match")
	}

}
