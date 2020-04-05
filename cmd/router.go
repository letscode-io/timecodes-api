package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func handleRequests() {
	log.Println("Starting development server at http://127.0.0.1:8080/")
	router := mux.NewRouter().StrictSlash(true)

	handler := cors.Default().Handler(router)

	router.HandleFunc("/annotation", createAnnotation).Methods("POST")
	router.HandleFunc("/annotations/{videoId}", getAnnotations)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

func createAnnotation(w http.ResponseWriter, r *http.Request) {
	var annotation Annotation

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &annotation)
	err := db.Create(&annotation).Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(annotation)
	}
}

func getAnnotations(w http.ResponseWriter, r *http.Request) {
	var annotations []Annotation

	vars := mux.Vars(r)
	key := vars["videoId"]

	err := db.Where(&Annotation{VideoID: key}).Find(&annotations).Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(annotations)
	}
}
