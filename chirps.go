package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/justinjest/chirpy/internal/auth"
	"github.com/justinjest/chirpy/internal/database"
	jsonParser "github.com/justinjest/chirpy/internal/json"
)

type chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	uuid, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	par := params{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&par)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if uuid != par.UserId {
		w.WriteHeader(401)
		return
	}
	if len(par.Body) > 140 {
		log.Printf("Chirp over 140 characters")
		return
	}
	tmp := par.Body
	par.Body = jsonParser.CleanBody(tmp)
	dbOutput := database.CreateChirpParams{
		Body:   par.Body,
		UserID: par.UserId,
	}
	newChirp, err := cfg.database.CreateChirp(context.Background(), dbOutput)
	if err != nil {
		fmt.Printf("Chirp unable to be stored %v\n", err)
	}
	tmpChirp := chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}
	dat, err := json.Marshal(tmpChirp)
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

	var res []chirp
	items, err := cfg.database.GetChirps(context.Background())
	if err != nil {
		fmt.Printf("Error getting cfg.database to load, %v\n", err)
		w.WriteHeader(400)
		return
	}
	for _, item := range items {
		tmp := chirp{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Body:      item.Body,
			UserID:    item.UserID,
		}
		res = append(res, tmp)
	}
	out, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Error marshalling res %v\n", err)
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
	w.Write(out)
}

func (cfg *apiConfig) getOneChirp(w http.ResponseWriter, r *http.Request) {
	pathString := r.PathValue("chirpID")
	pathUUID, err := uuid.Parse(pathString)
	fmt.Printf("%v, %v\n", pathString, pathUUID)
	fmt.Printf("%v\n", r)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		w.WriteHeader(500)
		return
	}
	chirpResp, err := cfg.database.GetOneChirp(context.Background(), pathUUID)
	if err != nil {
		fmt.Printf("Chrip not found %v\n", err)
		w.WriteHeader(404)
		return
	}
	res := chirp{
		ID:        chirpResp.ID,
		CreatedAt: chirpResp.CreatedAt,
		UpdatedAt: chirpResp.UpdatedAt,
		Body:      chirpResp.Body,
		UserID:    chirpResp.UserID,
	}
	out, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Error marshalling res %v\n", err)
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
	w.Write(out)
}
