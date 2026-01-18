package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateTitleRequest struct {
	VideoID string `json:"videoId"`
	Title   string `json:"title"`
}

func UpdateVideoTitle(c *gin.Context) {
	accessToken, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	token := accessToken.(string)

	var body UpdateTitleRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if body.VideoID == "" || body.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "videoId and title required"})
		return
	}

	// 1. Fetch existing snippet (required by YouTube)
	snippet, err := getVideoSnippet(body.VideoID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. Update title only
	snippet.Title = body.Title

	// 3. Push update
	if err := updateVideoSnippet(body.VideoID, snippet, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "title updated"})
}
