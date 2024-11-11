package auth

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := []byte("HelloWorld")
	hash, err := bcrypt.GenerateFromPassword(password, 1)
	if err != nil {
		log.Fatalf("Error with bcrypt %v\n", err)
	}
	err = bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		log.Fatalf("Error with comparison wanted %v, got %v\n", hash, err)
	}
}

func TestJWSAuth(t *testing.T) {
	new, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Error with uuid.NewUUID %v\n", err)
	}
	jwt, err := MakeJWT(new, "Hello", time.Hour*1)
	if err != nil {
		t.Errorf("Error with MakeJWT %v\n", err)
	}
	valid, err := ValidateJWT(jwt, "Hello")
	if err != nil {
		t.Errorf("Error with validateJWT %v\n", err)
	}
	if valid != new {
		t.Errorf("Error with validate got %v expected %v\n", valid, new)
	}
}

func TestJWSAuthExpired(t *testing.T) {
	new, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Error with uuid.NewUUID %v\n", err)
	}
	jwt, err := MakeJWT(new, "Hello", time.Millisecond*10)
	if err != nil {
		t.Errorf("Error with MakeJWT %v\n", err)
	}
	time.Sleep(time.Millisecond * 11)
	_, err = ValidateJWT(jwt, "Hello")
	if err == nil {
		t.Errorf("Error with expirationTime %v\n", err)
	}
}
func TestJWSAuthBadKeys(t *testing.T) {
	new, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Error with uuid.NewUUID %v\n", err)
	}
	jwt, err := MakeJWT(new, "Hello", time.Millisecond*10)
	if err != nil {
		t.Errorf("Error with MakeJWT %v\n", err)
	}
	time.Sleep(time.Millisecond * 11)
	val, err := ValidateJWT(jwt, "World")
	if err == nil {
		t.Errorf("same output with wrong keys %v, %v\n", jwt, val)
	}
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "Bearer abc123")
	_, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("Get Bearer Token Broke %v\n", err)
	}
}
