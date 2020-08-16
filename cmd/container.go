package main

import (
	youtubeapi "timecodes/cmd/youtube_api"
)

// Container represents DI container
type Container struct {
	TimecodeLikeRepository TimecodeLikeRepository
	TimecodeRepository     TimecodeRepository
	UserRepository         UserRepository

	YoutubeAPI youtubeapi.IService
}
