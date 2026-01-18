package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type CommentSnippet struct {
	AuthorDisplayName     string `json:"authorDisplayName"`
	AuthorProfileImageUrl string `json:"authorProfileImageUrl"`
	AuthorChannelUrl      string `json:"authorChannelUrl"`
	AuthorChannelId       struct {
		Value string `json:"value"`
	} `json:"authorChannelId"`
	TextOriginal string `json:"textOriginal"`
}

type TopLevelComment struct {
	Id      string         `json:"id"`
	Snippet CommentSnippet `json:"snippet"`
}

type ReplyComment struct {
	Id      string         `json:"id"`
	Snippet CommentSnippet `json:"snippet"`
}

type CommentThreadItem struct {
	Id      string `json:"id"`
	Snippet struct {
		ChannelId       string          `json:"channelId"`
		TopLevelComment TopLevelComment `json:"topLevelComment"`
	} `json:"snippet"`
	Replies struct {
		Comments []ReplyComment `json:"comments"`
	} `json:"replies"`
}

type YTCommentThreadResponse struct {
	NextPageToken string              `json:"nextPageToken"`
	Items         []CommentThreadItem `json:"items"`
}

func GetCommentThread(c *gin.Context) {
	accessToken, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Access token missing",
		})
		return
	}

	token, ok := accessToken.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid access token type",
		})
		return
	}

	videoId := c.Query("videoId")
	pageToken := c.Query("pageToken")

	if videoId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "videoId unavailable",
		})
		return
	}

	ytres, err := fetchComments(videoId, pageToken, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cant fetch comments",
		})
		return
	}

	fmt.Printf("Comments : %s\n", ytres.Items[0].Snippet.TopLevelComment.Snippet.TextOriginal)

	c.JSON(200, ytres)
}

func fetchComments(videoId string, pageToken string, token string) (*YTCommentThreadResponse, error) {
	reqURL, _ := url.Parse("https://www.googleapis.com/youtube/v3/commentThreads")

	q := reqURL.Query()
	q.Set("part", "snippet,replies")
	q.Set("videoId", videoId)

	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}

	reqURL.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("commentThreads error: %s", body)
	}

	var ytRes YTCommentThreadResponse
	if err := json.NewDecoder(res.Body).Decode(&ytRes); err != nil {
		return nil, err
	}

	return &ytRes, nil
}
