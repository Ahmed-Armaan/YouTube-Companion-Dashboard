package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type UserVideoListResponse struct {
	NextPageToken string  `json:"nextPageToken"`
	Videos        []Video `json:"videos"`
}

type YoutubeVideoResponse struct {
	Items []VideoItem `json:"items"`
}

type Video struct {
	ID            string `json:"videoId"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	PublishedAt   string `json:"publishedAt"`
	Thumbnail     string `json:"thumbnail"`
	Duration      string `json:"duration"`
	ViewCount     string `json:"viewCount"`
	LikesCount    string `json:"likesCount"`
	DisLikesCount string `json:"disLikesCount"`
	EmbeddedHTML  string `json:"embeddedhtml"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type YoutubeChannelResponse struct {
	Items []struct {
		ContentDetails struct {
			RelatedPlaylists struct {
				Uploads string `json:"uploads"`
			} `json:"relatedPlaylists"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type YoutubePlaylistItemsResponse struct {
	NextPageToken string         `json:"nextPageToken"`
	Items         []PlaylistItem `json:"items"`
}

type PlaylistItem struct {
	Snippet PlaylistItemSnippet `json:"snippet"`
}

type PlaylistItemSnippet struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	PublishedAt string               `json:"publishedAt"`
	ResourceID  ResourceID           `json:"resourceId"`
	Thumbnails  map[string]Thumbnail `json:"thumbnails"`
}

type ResourceID struct {
	VideoID string `json:"videoId"`
}

type VideoItem struct {
	ContentDetails ContentDetails `json:"contentDetails"`
	Statistics     Statistics     `json:"statistics"`
	Player         Player         `json:"player"`
}

type ContentDetails struct {
	Duration string `json:"duration"`
}

type Statistics struct {
	ViewCount    string `json:"viewCount"`
	LikeCount    string `json:"likeCount"`
	DislikeCount string `json:"dislikeCount,omitempty"`
}

type Player struct {
	EmbedHTML string `json:"embedHtml"`
}

//func Me(c *gin.Context) {
//	_, exists := c.Get("accessToken")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{
//			"error": "Not authenticated",
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"authenticated": true,
//	})
//}

func GetChannel(c *gin.Context) {
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

	reqURL, _ := url.Parse("https://www.googleapis.com/youtube/v3/channels")

	q := reqURL.Query()
	q.Set("part", "contentDetails")
	q.Set("mine", "true")
	reqURL.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "YouTube request failed"})
		return
	}
	defer res.Body.Close()

	var channelRes YoutubeChannelResponse
	if err := json.NewDecoder(res.Body).Decode(&channelRes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid YouTube response"})
		return
	}

	if len(channelRes.Items) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No channel found"})
		return
	}

	uploadPlaylistID := channelRes.Items[0].ContentDetails.RelatedPlaylists.Uploads

	nextPageToken := c.Query("pageToken")
	videoList, err := getVideosList(token, uploadPlaylistID, nextPageToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videoList)
}

func getVideosList(token string, uploadPlaylistID string, pageToken string) (*UserVideoListResponse, error) {
	reqURL, _ := url.Parse("https://www.googleapis.com/youtube/v3/playlistItems")

	q := reqURL.Query()
	q.Set("part", "snippet")
	q.Set("playlistId", uploadPlaylistID)
	q.Set("maxResults", "50")
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
		return nil, fmt.Errorf("youtube error: %s", body)
	}

	var ytRes YoutubePlaylistItemsResponse
	if err := json.NewDecoder(res.Body).Decode(&ytRes); err != nil {
		return nil, err
	}

	videos := make([]Video, 0, len(ytRes.Items))
	for _, item := range ytRes.Items {
		currVideoDetails, err := getVideo(item.Snippet.ResourceID.VideoID, token)
		if err != nil {
			continue
		}

		videos = append(videos, Video{
			ID:           item.Snippet.ResourceID.VideoID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  item.Snippet.PublishedAt,
			Thumbnail:    getThumbnail(item.Snippet.Thumbnails),
			Duration:     currVideoDetails.ContentDetails.Duration,
			ViewCount:    currVideoDetails.Statistics.ViewCount,
			LikesCount:   currVideoDetails.Statistics.LikeCount,
			EmbeddedHTML: currVideoDetails.Player.EmbedHTML,
		})
	}

	fmt.Printf("VideoData : \n%v\n", videos)

	return &UserVideoListResponse{
		NextPageToken: ytRes.NextPageToken,
		Videos:        videos,
	}, nil
}

func getThumbnail(t map[string]Thumbnail) string {
	if thumb, ok := t["maxres"]; ok {
		return thumb.URL
	}
	if thumb, ok := t["high"]; ok {
		return thumb.URL
	}
	if thumb, ok := t["medium"]; ok {
		return thumb.URL
	}
	if thumb, ok := t["default"]; ok {
		return thumb.URL
	}
	return ""
}

func getVideo(videoId string, token string) (VideoItem, error) {
	var resBody YoutubeVideoResponse

	reqURL, _ := url.Parse("https://www.googleapis.com/youtube/v3/videos")
	q := reqURL.Query()
	q.Set("part", "contentDetails,statistics,player")
	q.Set("id", videoId)
	reqURL.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return VideoItem{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return VideoItem{}, fmt.Errorf("videos.list failed")
	}

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return VideoItem{}, err
	}

	if len(resBody.Items) == 0 {
		return VideoItem{}, fmt.Errorf("no video data")
	}

	return resBody.Items[0], nil
}
