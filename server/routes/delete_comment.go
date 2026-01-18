package routes

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteComment(c *gin.Context) {
	commentId := c.Query("commentId")
	if commentId == "" {
		c.JSON(400, gin.H{"error": "commentId required"})
		return
	}

	tokenAny, exists := c.Get("accessToken")
	if !exists {
		c.JSON(401, gin.H{"error": "not authenticated"})
		return
	}
	token := tokenAny.(string)

	reqURL := "https://www.googleapis.com/youtube/v3/comments?id=" + commentId
	req, _ := http.NewRequest(http.MethodDelete, reqURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(res.Body)
		c.JSON(500, gin.H{"error": string(b)})
		return
	}

	c.JSON(200, gin.H{"status": "comment deleted"})
}
