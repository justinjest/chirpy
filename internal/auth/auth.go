package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	uuidString := userID.String()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   uuidString,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := t.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return str, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// First, verify the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the key as []byte
		return []byte(tokenSecret), nil
	}
	var claimStruct jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claimStruct, keyFunc)
	if err != nil {
		return uuid.Nil, err
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		u3, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, err
		}
		return u3, nil
	}
	return uuid.Nil, errors.New("Error")
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no authorization")
	}
	trim := strings.TrimPrefix(auth, "Bearer ")
	return trim, nil
}

func MakeRefreshToken() (string, error) {
	c := 32
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	return hex.EncodeToString(b), nil
}
