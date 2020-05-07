package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	timecodeParser "timecodes/cmd/timecode_parser"
)

type TimecodeJSON struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	LikesCount  int    `json:"likesCount"`
	LikedByMe   bool   `json:"likedByMe"`
	Seconds     int    `json:"seconds"`
	VideoID     string `json:"videoId"`
}

// GET /timecodes
func handleGetTimecodes(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
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
			timecodeJSONCollection = append(timecodeJSONCollection, serializeTimecode(timecode, currentUser))
		}

		json.NewEncoder(w).Encode(timecodeJSONCollection)
	}
}

// POST /timecodes
func handleCreateTimecode(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	timecode := &Timecode{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, timecode)
	err := db.Create(timecode).Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(serializeTimecode(timecode, currentUser))
	}
}

func serializeTimecode(timecode *Timecode, currentUser *User) (timecodeJSON *TimecodeJSON) {
	var likedByMe bool
	if currentUser != nil {
		likedByMe = getLikedByMe(timecode.Likes, currentUser.ID)
	}

	return &TimecodeJSON{
		ID:          timecode.ID,
		Description: timecode.Description,
		LikesCount:  len(timecode.Likes),
		LikedByMe:   likedByMe,
		Seconds:     timecode.Seconds,
		VideoID:     timecode.VideoID,
	}
}

func getLikedByMe(likes []TimecodeLike, userID uint) bool {
	for _, like := range likes {
		if like.UserID == userID {
			return true
		}
	}

	return false
}

func parseDescriptionAndCreateAnnotations(videoId string) {
	description := youtubeService.FetchVideoDescription(videoId)
	parsedCodes := timecodeParser.Parse(description)

	createTimecodes(parsedCodes, videoId)
}

func parseCommentsAndCreateAnnotations(videoId string) {
	var parsedCodes []timecodeParser.ParsedTimeCode

	comments, err := youtubeService.FetchVideoComments(videoId)
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
