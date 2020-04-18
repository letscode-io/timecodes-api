package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// GET /annotations
func getAnnotations(w http.ResponseWriter, r *http.Request) {
	annotations := &[]Annotation{}

	vars := mux.Vars(r)
	videoId := vars["videoId"]

	err := db.Order("seconds asc").Where(&Annotation{VideoID: videoId}).Find(annotations).Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(annotations)
	}
}

// POST /annotations
func createAnnotation(w http.ResponseWriter, r *http.Request) {
	annotation := &Annotation{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, annotation)
	err := db.Create(annotation).Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(annotation)
	}
}
