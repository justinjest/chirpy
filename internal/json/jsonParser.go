package jsonParser

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	type parameteres struct {
		Body string `body:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameteres{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	var errorMsg string
	if len(params.Body) > 140 {
		log.Printf("Chirp over 140 characters")
		errorMsg = "Chirp too long"
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
		Valid       bool   `json:"valid"`
	}
	output := returnVals{}
	if errorMsg != "" {
		output.CleanedBody = ""
		output.Valid = false
	} else {
		output.CleanedBody = cleanBody(params.Body)
		output.Valid = true
	}
	dat, err := json.Marshal(output)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if output.Valid {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	w.Write(dat)
}

func cleanBody(s string) string {
	const censorItem = "****"
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	split := strings.Split(s, " ")
	var tmp []string
	for _, word := range split {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			tmp = append(tmp, censorItem)
		} else {
			tmp = append(tmp, word)
		}
	}
	return strings.Join(tmp, " ")
}
