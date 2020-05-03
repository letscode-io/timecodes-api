package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func handleCreateTimecodeLike(w http.ResponseWriter, r *http.Request) {
	like := &TimecodeLike{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, like)
	err := db.Create(like).Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(like)
	}
}
