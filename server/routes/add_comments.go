package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddCommentRequest struct {
	VideoID string `json:"videoId"`
	Text    string `json:"text"`
}

func AddComment(c *gin.Context) {
	tokenAny, exists := c.Get("accessToken")
	if !exists {
		c.JSON(401, gin.H{"error": "not authenticated"})
		return
	}
	token := tokenAny.(string)

	var body AddCommentRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	payload := map[string]any{
		"snippet": map[string]any{
			"videoId": body.VideoID,
			"topLevelComment": map[string]any{
				"snippet": map[string]any{
					"textOriginal": body.Text,
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest(
		http.MethodPost,
		"https://www.googleapis.com/youtube/v3/commentThreads?part=snippet",
		bytes.NewReader(jsonBody),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		c.JSON(500, gin.H{"error": string(b)})
		return
	}

	c.JSON(200, gin.H{"status": "comment added"})
}
