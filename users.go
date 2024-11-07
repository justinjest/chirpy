package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
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
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameteres{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding paramaters %v", err)
		return
	}
	usr, err := cfg.database.CreateUser(context.Background(), params.Email)
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
