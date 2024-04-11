package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
)

const (
	API_PORT = 8081
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.FormValue("cep")

	if len(cep) != 8 || reflect.TypeOf(cep).String() != "string" {
		log.Printf("invalid zipcode")
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	fmt.Println("orchestrator: " + cep)
}

func main() {
	http.HandleFunc("/weather", WeatherHandler)
	log.Printf("Listening on %d", API_PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", API_PORT), nil)
}
