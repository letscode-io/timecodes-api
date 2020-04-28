package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	timecodeParser "youannoapi/cmd/timecode_parser"

	"github.com/gorilla/mux"
)

// GET /timecodes
func getTimecodes(w http.ResponseWriter, r *http.Request) {
	timecodes := &[]Timecode{}

	vars := mux.Vars(r)
	videoId := vars["videoId"]

	err := db.Order("seconds asc").Where(&Timecode{VideoID: videoId}).Find(timecodes).Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		if len(*timecodes) == 0 {
			go func() {
				parseDescriptionAndCreateAnnotations(videoId)
				parseCommentsAndCreateAnnotations(videoId)
			}()
		}

		json.NewEncoder(w).Encode(timecodes)
	}
}

// POST /timecodes
func createTimecode(w http.ResponseWriter, r *http.Request) {
	timecode := &Timecode{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, timecode)
	err := db.Create(timecode).Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(timecode)
	}
}

func parseDescriptionAndCreateAnnotations(videoId string) {
	description := getVideoDescription(videoId)
	parsedCodes := timecodeParser.Parse(description)

	for _, code := range parsedCodes {
		timecode := &Timecode{Seconds: code.Seconds, VideoID: videoId, Description: code.Description}
		err := db.Create(timecode).Error
		if err != nil {
			log.Println(err)
		}
	}
}

func parseCommentsAndCreateAnnotations(videoId string) {
	var parsedCodes []timecodeParser.ParsedTimeCode

	comments, err := fetchVideoComments(videoId)
	if err != nil {
		log.Println(err)
	}

	for _, comment := range comments {
		timeCodes := timecodeParser.Parse(comment.Snippet.TopLevelComment.Snippet.TextOriginal)

		parsedCodes = append(parsedCodes, timeCodes...)
	}

	for _, code := range parsedCodes {
		annotation := &Timecode{Seconds: code.Seconds, VideoID: videoId, Description: code.Description}
		err := db.Create(annotation).Error
		if err != nil {
			log.Println(err)
		}
	}
}
