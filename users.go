package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/justinjest/chirpy/internal/auth"
	"github.com/justinjest/chirpy/internal/database"
)

func (cfg *apiConfig) DropUsers(w http.ResponseWriter, req *http.Request) {
	if cfg.PLATFORM != "dev" {
		w.WriteHeader(403)
		return
	}
	cfg.database.DeleteUsers(context.Background())
	w.WriteHeader(200)
}

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, req *http.Request) {
	type parameteres struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameteres{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding paramaters %v", err)
		return
	}
	params.Password, err = auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("errror generating hash %v", err)
		return
	}
	usr, err := cfg.database.CreateUser(context.Background(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: params.Password,
	})
	if err != nil {
		log.Printf("error creating user %v", err)
		return
	}
	w.WriteHeader(201)
	res := User{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshiling user %v", err)
		return
	}
	w.Write(data)
}

func (apiCfg *apiConfig) userLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding paramaters %v", err)
		return
	}
	hash, err := apiCfg.database.GetUserHash(context.Background(), params.Email)
	if err != nil {
		log.Printf("error retriving hash %v", err)
		w.WriteHeader(500)
		return
	}
	err = auth.CheckPasswordHash(params.Password, hash)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("Incorrect email and password"))
		return
	}
	usr, err := apiCfg.database.GetUserFromEmail(context.Background(), params.Email)
	if err != nil {
		log.Printf("error creating user %v", err)
		return
	}
	token, err := auth.MakeJWT(usr.ID, apiCfg.secret)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		fmt.Printf("error creating refresh token %v\n", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: usr.ID,
	}
	tkn, err := apiCfg.database.CreateRefreshToken(context.Background(), refreshTokenParams)
	if err != nil {
		fmt.Print(tkn)
		fmt.Print("Unable to generate apiCFG database, ", err)
		w.WriteHeader(500)
		return
	}
	res := User{
		ID:           usr.ID,
		CreatedAt:    usr.CreatedAt,
		UpdatedAt:    usr.UpdatedAt,
		Email:        usr.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}
	fmt.Printf("userID: %v\n", res.ID)
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshiling user %v", err)
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}

func (apiCfg *apiConfig) refreshUser(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Print("error, ", err)
		w.WriteHeader(401)
		return
	}
	usr, err := apiCfg.database.GetUserFromRefreshToken(context.Background(), token)
	if err != nil {
		fmt.Print(token)
		log.Print("error getting user from token db", err)
		w.WriteHeader(401)
		return
	}
	auth, err := auth.MakeJWT(usr.ID, apiCfg.secret)
	if err != nil {
		log.Print("error making JWT, ", err)
		w.WriteHeader(401)
		return
	}
	type response struct {
		Token string `json:"token"`
	}
	res := response{
		Token: auth,
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshiling user %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
