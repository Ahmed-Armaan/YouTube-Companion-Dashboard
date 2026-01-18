package database

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"yt_dashboard.com/utils"
)

func InsertUser(googleUserId, name, email string) (uuid.UUID, error) {
	var user User

	err := DB.Where("google_user_id = ?", googleUserId).First(&user).Error
	if err == nil {
		return user.ID, nil
	}
	if err != gorm.ErrRecordNotFound {
		return uuid.Nil, err
	}

	user = User{
		GoogleUserID: googleUserId,
		Name:         name,
		Email:        email,
		CreatedAt:    time.Now(),
	}

	if err := DB.Create(&user).Error; err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func InsertToken(userID uuid.UUID, refreshToken string) error {
	encToken, err := utils.Encrypt(refreshToken)
	if err != nil {
		return err
	}

	var token Token
	err = DB.First(&token, "user_id = ?", userID).Error

	if err == nil {
		token.RefreshTokenEnc = encToken
		token.Revoked = false
		token.CreatedAt = time.Now()
		return DB.Save(&token).Error
	}

	if err != gorm.ErrRecordNotFound {
		return err
	}

	token = Token{
		UserID:          userID,
		RefreshTokenEnc: encToken,
		Revoked:         false,
		CreatedAt:       time.Now(),
	}

	return DB.Create(&token).Error
}

func GetToken(googleUserId string) (string, error) {
	// var token Token
	// err := DB.First(&token, "user_id = ? AND revoked = false", userID).Error
	//
	//	if err != nil {
	//		return "", errors.New("no valid refresh token")
	//	}
	// return utils.Decrypt(token.RefreshTokenEnc)

	var user User
	err := DB.First(&user, "google_user_id = ?", googleUserId).Error
	if err != nil {
		return "", errors.New("No valid refresh token present")
	}

	var token Token
	err = DB.First(&token, "user_id = ? AND revoked = false", user.ID).Error
	if err != nil {
		return "", errors.New("No valid refresh token present")
	}

	return utils.Decrypt(token.RefreshTokenEnc)
}
