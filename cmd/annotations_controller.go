package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	timecodeParser "youannoapi/cmd/timecode_parser"

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
		if len(*annotations) == 0 {
			go func() {
				parseDescriptionAndCreateAnnotations(videoId)
				parseCommentsAndCreateAnnotations(videoId)
			}()
		}

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

func parseDescriptionAndCreateAnnotations(videoId string) {
	description := getVideoDescription(videoId)
	parsedCodes := timecodeParser.Parse(description)

	for _, code := range parsedCodes {
		annotation := &Annotation{Seconds: code.Seconds, VideoID: videoId, Text: code.Description}
		err := db.Create(annotation).Error
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
		annotation := &Annotation{Seconds: code.Seconds, VideoID: videoId, Text: code.Description}
		err := db.Create(annotation).Error
		if err != nil {
			log.Println(err)
		}
	}
}
