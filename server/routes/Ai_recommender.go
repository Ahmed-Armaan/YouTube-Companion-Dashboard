package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type TitleAIRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TitleAIResponse struct {
	Suggestions []string `json:"suggestions"`
}

func SuggestTitles(c *gin.Context) {
	var body TitleAIRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OPENAI_API_KEY not set"})
		return
	}

	prompt := "Improve the following YouTube title. Return exactly 3 concise, catchy alternatives.\n\n" +
		"Title: " + body.Title + "\n" +
		"Description: " + body.Description

	payload := map[string]any{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode request"})
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.openai.com/v1/chat/completions",
		bytes.NewReader(b),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		buf.ReadFrom(res.Body)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ai request failed",
			"body":  buf.String(),
		})
		return
	}

	var aiRes struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(res.Body).Decode(&aiRes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid ai response"})
		return
	}

	if len(aiRes.Choices) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ai returned no choices"})
		return
	}

	raw := aiRes.Choices[0].Message.Content

	lines := []string{}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "-")
		line = strings.TrimPrefix(line, "•")
		line = strings.TrimSpace(line)

		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) > 3 {
		lines = lines[:3]
	}

	c.JSON(http.StatusOK, TitleAIResponse{
		Suggestions: lines,
	})
}
