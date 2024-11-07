package jsonParser

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameteres struct {
	Body  string `json:"body"`
	Valid bool   `json:"valid"`
}
type ErrorVal struct {
	Errormsg string `json:"error"`
}

func ChirpValidator(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameteres{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if len(params.Body) > 140 {
		log.Printf("Chirp over 140 characters")
		res := ErrorVal{
			Errormsg: "Chirp over 140 characters",
		}
		errorWriter(res, w)
		return
	}

	output := parameteres{
		Body:  CleanBody(params.Body),
		Valid: true,
	}
	dat, err := json.Marshal(output)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		res := ErrorVal{
			Errormsg: "Something went wrong",
		}
		errorWriter(res, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func errorWriter(errorJson ErrorVal, w http.ResponseWriter) {
	output := errorJson
	dat, err := json.Marshal(output)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(dat)
}

func CleanBody(s string) string {
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
