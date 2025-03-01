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
		Body string `json:"body"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		fmt.Printf("Error getting bearer token %v\n", token)
		return
	}
	uuid, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		w.WriteHeader(401)
		fmt.Printf("Unable to validate token, %v\n", err)
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
	if len(par.Body) > 140 {
		log.Printf("Chirp over 140 characters")
		return
	}
	tmp := par.Body
	par.Body = jsonParser.CleanBody(tmp)
	dbOutput := database.CreateChirpParams{
		Body:   par.Body,
		UserID: uuid,
	}
	newChirp, err := cfg.database.CreateChirp(context.Background(), dbOutput)
	if err != nil {
		fmt.Printf("Chirp unable to be stored %v\n", err)
		w.WriteHeader(500)
		return
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
func (cfg *apiConfig) getAllChirpsASC(w http.ResponseWriter, r *http.Request) {
	items, err := cfg.database.GetChirpsASC(context.Background())
	if err != nil {
		fmt.Printf("Error getting cfg.database to load, %v\n", err)
		w.WriteHeader(400)
		return
	}
	cfg.writeChirps(w, r, items)
}
func (cfg *apiConfig) getAllChirpsByIDasc(w http.ResponseWriter, r *http.Request, user_id uuid.UUID) {
	items, err := cfg.database.GetChirpsByIDASC(context.Background(), user_id)
	if err != nil {
		fmt.Printf("Error getting cfg.database to load, %v\n", err)
		w.WriteHeader(400)
		return
	}
	cfg.writeChirps(w, r, items)
}
func (cfg *apiConfig) getAllChirpsdesc(w http.ResponseWriter, r *http.Request) {
	items, err := cfg.database.GetChirpsDESC(context.Background())
	if err != nil {
		fmt.Printf("Error getting cfg.database to load, %v\n", err)
		w.WriteHeader(400)
		return
	}
	cfg.writeChirps(w, r, items)
}
func (cfg *apiConfig) getAllChirpsByIDdesc(w http.ResponseWriter, r *http.Request, user_id uuid.UUID) {
	items, err := cfg.database.GetChirpsByIDDESC(context.Background(), user_id)
	if err != nil {
		fmt.Printf("Error getting cfg.database to load, %v\n", err)
		w.WriteHeader(400)
		return
	}
	cfg.writeChirps(w, r, items)
}
func (cfg *apiConfig) writeChirps(w http.ResponseWriter, r *http.Request, items []database.Chirp) {
	var res []chirp
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

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	ord := r.URL.Query().Get("sort")
	fmt.Printf("%v\n", ord)
	if ord == "" || ord == "asc" {
		if s == "" {
			cfg.getAllChirpsASC(w, r)
		}
		user_id, err := uuid.Parse(s)
		if err != nil {
			log.Printf("Error parsing uuid %v\n", err)
			return
		}
		cfg.getAllChirpsByIDasc(w, r, user_id)
	}
	if ord == "desc" {
		if s == "" {
			cfg.getAllChirpsdesc(w, r)
		}
		user_id, err := uuid.Parse(s)
		if err != nil {
			log.Printf("Error parsing uuid %v\n", err)
			return
		}
		cfg.getAllChirpsByIDdesc(w, r, user_id)
	}
}
func (cfg *apiConfig) getOneChirp(w http.ResponseWriter, r *http.Request) {
	pathString := r.PathValue("chirpID")
	pathUUID, err := uuid.Parse(pathString)
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

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	pathString := r.PathValue("chirpID")
	pathUUID, err := uuid.Parse(pathString)
	if err != nil {
		fmt.Printf("Error parsing path: %v\n", err)
		w.WriteHeader(500)
		return
	}
	chirpResp, err := cfg.database.GetOneChirp(context.Background(), pathUUID)
	if err != nil {
		fmt.Printf("Chrip not found %v\n", err)
		w.WriteHeader(404)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Print("error getting bearer token, ", err)
		w.WriteHeader(401)
		return
	}
	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error validating jwt %v\n", err)
		w.WriteHeader(401)
		return
	}
	if chirpResp.UserID != id {
		log.Printf("User %v attempted to access %v", chirpResp.UserID, id)
		w.WriteHeader(403)
		return
	}
	err = cfg.database.DeleteChirp(context.Background(), chirpResp.ID)
	if err != nil {
		log.Printf("Unable to find chirp %v\n", err)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}
