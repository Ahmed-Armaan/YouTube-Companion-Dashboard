package database

import (
	"github.com/lib/pq"
	"time"
)

func InsertNote(note *Note) error {
	return DB.Create(note).Error
}

func GetNotes(videoId string, limit int, cursor *time.Time) ([]Note, error) {
	var notes []Note

	query := DB.Where("video_id = ?", videoId).Order("created_at DESC").Limit(limit)
	if cursor != nil {
		query = query.Where("created_at < ?", *cursor)
	}

	err := query.Find(&notes).Error
	return notes, err
}

func GetNotesByTags(videoId string, limit int, tags []string, cursor *time.Time) ([]Note, error) {
	var notes []Note

	query := DB.Where("video_id = ?", videoId).Where("tags @> ?", pq.StringArray(tags)).Order("created_at DESC").Limit(limit)
	if cursor != nil {
		query = query.Where("created_at < ?", *cursor)
	}

	err := query.Find(&notes).Error
	return notes, err
}
