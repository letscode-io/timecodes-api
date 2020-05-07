package youtubeapi

import (
	"context"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Service struct {
	client *youtube.Service
}

func New() *Service {
	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(os.Getenv("GOOGLE_API_KEY")))
	if err != nil {
		log.Println(err)
	}

	return &Service{client: youtubeService}
}

func (s *Service) FetchVideoDescription(videoId string) string {
	call := s.client.
		Videos.
		List("snippet").
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

func (s *Service) FetchVideoComments(videoId string) ([]*youtube.CommentThread, error) {
	call := s.client.
		CommentThreads.
		List("snippet").
		VideoId(videoId).
		Order("relevance").
		MaxResults(100)

	response, err := call.Do()
	if err != nil {
		log.Println(err)

		return nil, err
	}

	return response.Items, nil
}
