package container

import (
	timecodeLikes "timecodes/pkg/timecode_likes"
	timecodes "timecodes/pkg/timecodes"
	"timecodes/pkg/users"
	youtubeapi "timecodes/pkg/youtube_api"
)

// Container represents DI container
type Container struct {
	TimecodeLikeRepository timecodeLikes.TimecodeLikeRepository
	TimecodeRepository     timecodes.TimecodeRepository
	UserRepository         users.UserRepository

	YoutubeAPI youtubeapi.IService
}
