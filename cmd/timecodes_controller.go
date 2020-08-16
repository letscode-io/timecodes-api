package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	timecodeParser "timecodes/pkg/timecode_parser"

	"github.com/gorilla/mux"
)

// TimecodeRequest represents timecode request parameters
type TimecodeRequest struct {
	Description string `json:"description"`
	RawSeconds  string `json:"seconds"`
	VideoID     string `json:"videoId"`
}

// TimecodeJSON represents custom timecode response
type TimecodeJSON struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	LikesCount  int    `json:"likesCount"`
	LikedByMe   bool   `json:"likedByMe,omitempty"`
	Seconds     int    `json:"seconds"`
	VideoID     string `json:"videoId"`
	UserID      uint   `json:"userId,omitempty"`
}

// GET /timecodes/{videoId}
func handleGetTimecodes(c *Container, w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	videoID := mux.Vars(r)["videoId"]

	timecodes := c.TimecodeRepository.FindByVideoID(videoID)

	if len(*timecodes) == 0 {
		go parseVideoContentAndCreateTimecodes(c, videoID)
	}

	timecodeJSONCollection := make([]*TimecodeJSON, 0)
	for _, timecode := range *timecodes {
		timecodeJSONCollection = append(timecodeJSONCollection, serializeTimecode(timecode, currentUser))
	}

	json.NewEncoder(w).Encode(timecodeJSONCollection)
}

// POST /timecodes
func handleCreateTimecode(c *Container, w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)

	reqBody, _ := ioutil.ReadAll(r.Body)
	timecodeRequest := &TimecodeRequest{}
	err := json.Unmarshal(reqBody, timecodeRequest)
	if err != nil {
		log.Println(err)

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	timecode := &Timecode{
		Description: timecodeRequest.Description,
		Seconds:     timecodeParser.ParseSeconds(timecodeRequest.RawSeconds),
		VideoID:     timecodeRequest.VideoID,
		UserID:      currentUser.ID,
	}

	_, err = c.TimecodeRepository.Create(timecode)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(serializeTimecode(timecode, currentUser))
	}
}

func serializeTimecode(timecode *Timecode, currentUser *User) (timecodeJSON *TimecodeJSON) {
	var likedByMe bool
	var userID uint

	if currentUser != nil {
		likedByMe = getLikedByMe(timecode.Likes, currentUser.ID)
		userID = timecode.UserID
	}

	return &TimecodeJSON{
		ID:          timecode.ID,
		Description: timecode.Description,
		LikesCount:  len(timecode.Likes),
		LikedByMe:   likedByMe,
		Seconds:     timecode.Seconds,
		VideoID:     timecode.VideoID,
		UserID:      userID,
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

func parseVideoContentAndCreateTimecodes(c *Container, videoID string) {
	description := c.YoutubeAPI.FetchVideoDescription(videoID)
	parsedCodes := timecodeParser.Parse(description)

	comments := c.YoutubeAPI.FetchVideoComments(videoID)

	for _, comment := range comments {
		timeCodes := timecodeParser.Parse(comment)

		parsedCodes = append(parsedCodes, timeCodes...)
	}

	c.TimecodeRepository.CreateFromParsedCodes(parsedCodes, videoID)
}
