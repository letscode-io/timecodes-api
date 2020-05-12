package main

import (
	"encoding/json"
	"io/ioutil"
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
func handleGetTimecodes(c *Container, w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	videoID := mux.Vars(r)["videoId"]

	timecodes := c.TimecodeRepository.FindByVideoId(videoID)

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
	timecode := &Timecode{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, timecode)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	_, err = c.TimecodeRepository.Create(timecode)
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
