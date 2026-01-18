package utils

import (
	"sync"
)

var (
	tokenMap = make(map[string]string)
	mu       sync.RWMutex
)

func GetAccessTokenFromCache(userId string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()

	if token, ok := tokenMap[userId]; ok {
		return token, true
	} else {
		return "", false
	}
}

func InsertAccessToken(userId string, token string) {
	mu.Lock()
	tokenMap[userId] = token
	mu.Unlock()
}
