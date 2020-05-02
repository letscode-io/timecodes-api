package main

import (
	"log"

	"google.golang.org/api/youtube/v3"
)

func getVideoDescription(videoId string) string {
	call := youtubeService.Videos.List("snippet")
	call = call.Id(videoId)

	response, err := call.Do()
	if err != nil {
		log.Println(err)

		return ""
	}

	items := response.Items
	if len(items) == 0 {
		return ""
	}

	return items[0].Snippet.Description
}

func fetchVideoComments(videoId string) ([]*youtube.CommentThread, error) {
	call := youtubeService.CommentThreads.List("snippet")
	call = call.VideoId(videoId)
	call = call.Order("relevance")
	call = call.MaxResults(100)

	response, err := call.Do()
	if err != nil {
		log.Println(err)

		return nil, err
	}

	return response.Items, nil
}
