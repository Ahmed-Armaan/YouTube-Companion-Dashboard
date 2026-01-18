package database

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	GoogleUserID string    `gorm:"uniqueIndex;not null"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"not null"`
	CreatedAt    time.Time
}

type Token struct {
	UserID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	RefreshTokenEnc string    `gorm:"not null"`
	Revoked         bool      `gorm:"default:false"`
	CreatedAt       time.Time
}

type Note struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	VideoID   string    `gorm:"index"`
	Content   string
	Tags      pq.StringArray `gorm:"type:text[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
