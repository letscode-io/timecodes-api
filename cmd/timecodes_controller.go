package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	timecodeParser "youannoapi/cmd/timecode_parser"

	"github.com/gorilla/mux"
)

type TimecodeJSON struct {
	Description string `json:"description"`
	LikesCount  int    `json:"likesCount"`
	Seconds     int    `json:"seconds"`
	VideoID     string `json:"videoId"`
}

// GET /timecodes
func getTimecodes(w http.ResponseWriter, r *http.Request) {
	timecodes := &[]*Timecode{}

	videoId := mux.Vars(r)["videoId"]

	err := db.Order("seconds asc").
		Preload("Likes").
		Where(&Timecode{VideoID: videoId}).
		Find(timecodes).
		Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		if len(*timecodes) == 0 {
			go func() {
				parseDescriptionAndCreateAnnotations(videoId)
				parseCommentsAndCreateAnnotations(videoId)
			}()
		}

		timecodeJSONCollection := make([]*TimecodeJSON, 0)
		for _, timecode := range *timecodes {
			timecodeJSONCollection = append(timecodeJSONCollection, serializeTimecode(timecode))
		}

		json.NewEncoder(w).Encode(timecodeJSONCollection)
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
		json.NewEncoder(w).Encode(serializeTimecode(timecode))
	}
}

func serializeTimecode(timecode *Timecode) (timecodeJSON *TimecodeJSON) {
	return &TimecodeJSON{
		Description: timecode.Description,
		LikesCount:  len(timecode.Likes),
		Seconds:     timecode.Seconds,
		VideoID:     timecode.VideoID,
	}
}

func parseDescriptionAndCreateAnnotations(videoId string) {
	description := getVideoDescription(videoId)
	parsedCodes := timecodeParser.Parse(description)

	createTimecodes(parsedCodes, videoId)
}

func parseCommentsAndCreateAnnotations(videoId string) {
	var parsedCodes []timecodeParser.ParsedTimeCode

	comments, err := fetchVideoComments(videoId)
	if err != nil {
		log.Println(err)

		return
	}

	for _, comment := range comments {
		timeCodes := timecodeParser.Parse(comment.Snippet.TopLevelComment.Snippet.TextOriginal)

		parsedCodes = append(parsedCodes, timeCodes...)
	}

	createTimecodes(parsedCodes, videoId)
}

func createTimecodes(parsedTimecodes []timecodeParser.ParsedTimeCode, videoId string) {
	seen := make(map[string]struct{})

	for _, code := range parsedTimecodes {
		key := string(code.Seconds) + code.Description
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}

		annotation := &Timecode{Seconds: code.Seconds, VideoID: videoId, Description: code.Description}
		err := db.Create(annotation).Error
		if err != nil {
			log.Println(err)
		}
	}
}
