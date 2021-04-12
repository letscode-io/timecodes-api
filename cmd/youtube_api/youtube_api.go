package youtubeapi

import (
	"context"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const GOOGLE_API_KEY = "GOOGLE_API_KEY"

type IService interface {
	FetchVideoDescription(string) string
	FetchVideoComments(string) []string
}

type Service struct {
	IService

	client *youtube.Service
}

func New() (*Service, error) {
	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(os.Getenv(GOOGLE_API_KEY)))
	if err != nil {
		return nil, err
	}

	return &Service{client: youtubeService}, nil
}

func (s *Service) FetchVideoDescription(videoId string) string {
	call := s.client.
		Videos.
		List([]string{"snippet"}).
		Id(videoId)

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

func (s *Service) FetchVideoComments(videoId string) []string {
	textComments := make([]string, 0)

	call := s.client.
		CommentThreads.
		List([]string{"snippet"}).
		VideoId(videoId).
		Order("relevance").
		MaxResults(100)

	response, err := call.Do()
	if err != nil {
		log.Println(err)

		return textComments
	}

	for _, item := range response.Items {
		textComments = append(textComments, item.Snippet.TopLevelComment.Snippet.TextOriginal)
	}

	return textComments
}
