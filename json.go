package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type chirp struct {
	// these tags indicate how the keys in the JSON should be mapped to the struct fields
	// the struct fields must be exported (start with a capital letter) if you want them parsed
	Body string `json:"body"`
}

type resBody struct {
	Error string `json:"error"`
}

type resBody2 struct {
	Valid bool `json:"valid"`
}

func validateChirp(w http.ResponseWriter, req *http.Request) {
	// Prepare response Body JSONs
	genErr := resBody{Error: "Something went wrong"}
	jsonGenErr, err := json.Marshal(genErr)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	lengthErr := resBody{Error: "Chirp is too long"}
	jsonLengthErr, err := json.Marshal(lengthErr)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	valid := resBody2{Valid: true}
	jsonValid, err := json.Marshal(valid)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	// Decode JSON Request Body
	decoder := json.NewDecoder(req.Body)
	chirp := chirp{}
	err = decoder.Decode(&chirp)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		w.Write(jsonGenErr)
		return
	}
	// Check length of Chirp
	runes := []rune(chirp.Body)
	if len(runes) > 140 { // chirp too long
		log.Printf("Chirp is too long (%v characters past 140):\n%v\n", len(runes)-140, chirp.Body)
		w.WriteHeader(400)
		w.Write(jsonLengthErr)
		return
	}
	// All good
	w.WriteHeader(200)
	w.Write(jsonValid)
}
