package main

import (
	"log"
)

func getVideoDescription(videoId string) string {
	call := youtubeService.Videos.List("snippet")
	call = call.Id(videoId)

	response, err := call.Do()
	if err != nil {
		log.Println(err)
	}

	return response.Items[0].Snippet.Description
}
