package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"timecodes/pkg/container"
	m "timecodes/pkg/models"
	"timecodes/pkg/users"
)

// HandleCreateTimecodeLike POST /timecode_likes
func HandleCreateTimecodeLike(c *container.Container, w http.ResponseWriter, r *http.Request) {
	currentUser := users.GetCurrentUser(r)
	like := &m.TimecodeLike{}

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

// HandleDeleteTimecodeLike DELETE /timecode_likes
func HandleDeleteTimecodeLike(c *container.Container, w http.ResponseWriter, r *http.Request) {
	currentUser := users.GetCurrentUser(r)
	timecodeLike := &m.TimecodeLike{}

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
