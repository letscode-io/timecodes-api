package main

import (
	"github.com/jinzhu/gorm"
)

func runMigrations(db *gorm.DB) {
	applyTimecodesMigrations(db)
	applyUsersMigrations(db)
	applyTimecodeLikesMigrations(db)
}

func applyTimecodesMigrations(db *gorm.DB) {
	db.AutoMigrate(&Timecode{})
	db.Model(&Timecode{}).AddUniqueIndex(
		"idx_timecodes_seconds_text_video_id",
		"seconds", "description", "video_id",
	)
	db.Exec(`
		ALTER TABLE timecodes
		ADD CONSTRAINT description_min_length CHECK (length(description) >= 1);
	`)
}

func applyUsersMigrations(db *gorm.DB) {
	db.AutoMigrate(&User{})

	getAdminUser(db)
}

func applyTimecodeLikesMigrations(db *gorm.DB) {
	db.AutoMigrate(&TimecodeLike{})
	db.Model(&TimecodeLike{}).AddUniqueIndex(
		"idx_timecodes_likes_user_id_timecode_id_video_id",
		"user_id", "timecode_id",
	)
	db.Model(&TimecodeLike{}).AddForeignKey("timecode_id", "timecodes(id)", "RESTRICT", "RESTRICT")
	db.Model(&TimecodeLike{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
}
