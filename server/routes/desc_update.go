package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type UpdateDescriptionRequest struct {
	VideoID     string `json:"videoId"`
	Description string `json:"description"`
}

type VideoSnippet struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryID  string `json:"categoryId"`
}

func UpdateVideoDescription(c *gin.Context) {
	accessToken, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	token := accessToken.(string)

	var body UpdateDescriptionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if body.VideoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "videoId required"})
		return
	}

	snippet, err := getVideoSnippet(body.VideoID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	snippet.Description = body.Description

	if err := updateVideoSnippet(body.VideoID, snippet, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "description updated"})
}

func getVideoSnippet(videoId, token string) (VideoSnippet, error) {
	reqURL, _ := url.Parse("https://www.googleapis.com/youtube/v3/videos")
	q := reqURL.Query()
	q.Set("part", "snippet")
	q.Set("id", videoId)
	reqURL.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return VideoSnippet{}, err
	}
	defer res.Body.Close()

	var resp struct {
		Items []struct {
			Snippet VideoSnippet `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return VideoSnippet{}, err
	}

	if len(resp.Items) == 0 {
		return VideoSnippet{}, fmt.Errorf("video not found")
	}

	return resp.Items[0].Snippet, nil
}

func updateVideoSnippet(videoId string, snippet VideoSnippet, token string) error {
	body := map[string]any{
		"id":      videoId,
		"snippet": snippet,
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		http.MethodPut,
		"https://www.googleapis.com/youtube/v3/videos?part=snippet",
		bytes.NewReader(jsonBody),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("update failed: %s", b)
	}

	return nil
}
