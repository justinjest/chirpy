package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	jsonParser "github.com/justinjest/chirpy/internal/json"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	par := params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&par)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(par.Body) > 140 {
		log.Printf("Chirp over 140 characters")
		return
	}
	tmp := par.Body
	par.Body = jsonParser.CleanBody(tmp)
	dat, err := json.Marshal(par)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	return
}
