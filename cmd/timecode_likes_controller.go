package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// POST /timecode_likes
func handleCreateTimecodeLike(c *Container, w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	like := &TimecodeLike{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, like)
	if err != nil {
		log.Println(err)

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = c.TimecodeLikeRepository.Create(like, currentUser.ID)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(like)
	}
}

// DELETE /timecode_likes
func handleDeleteTimecodeLike(c *Container, w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	timecodeLike := &TimecodeLike{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, timecodeLike)
	if err != nil {
		log.Println(err)

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = c.TimecodeLikeRepository.Delete(timecodeLike, currentUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(timecodeLike)
	}
}
