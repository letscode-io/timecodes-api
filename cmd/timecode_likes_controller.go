package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func handleCreateTimecodeLike(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	like := &TimecodeLike{UserID: currentUser.ID}

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, like)
	if err != nil {
		log.Println(err)
		return
	}

	err = db.Create(like).Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(like)
	}
}

func handleDeleteTimecodeLike(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	likeParams := &TimecodeLike{}
	like := &TimecodeLike{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, likeParams)
	if err != nil {
		log.Println(err)
		return
	}

	err = db.Where(&TimecodeLike{UserID: currentUser.ID, TimecodeID: likeParams.TimecodeID}).First(like).Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	err = db.Unscoped().Delete(like).Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(like)
	}
}
