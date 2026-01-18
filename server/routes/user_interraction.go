package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	_, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
	})
}

func MyChannelId(c *gin.Context) {
	token := c.MustGet("accessToken").(string)

	req, _ := http.NewRequest(
		"GET",
		"https://www.googleapis.com/youtube/v3/channels?part=id&mine=true",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "youtube request failed"})
		return
	}
	defer res.Body.Close()

	var out struct {
		Items []struct {
			Id string `json:"id"`
		} `json:"items"`
	}

	json.NewDecoder(res.Body).Decode(&out)

	c.JSON(200, gin.H{
		"channelId": out.Items[0].Id,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie(
		"session",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"message": "logged out",
	})
}
