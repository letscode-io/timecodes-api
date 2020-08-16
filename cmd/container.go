package main

import (
	youtubeapi "timecodes/pkg/youtube_api"
)

// Container represents DI container
type Container struct {
	TimecodeLikeRepository TimecodeLikeRepository
	TimecodeRepository     TimecodeRepository
	UserRepository         UserRepository

	YoutubeAPI youtubeapi.IService
}
