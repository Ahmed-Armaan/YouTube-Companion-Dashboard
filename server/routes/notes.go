package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"yt_dashboard.com/database"
)

type CreateNoteRequest struct {
	VideoID string   `json:"videoId"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func CreateNote(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	userID := userIDAny.(uuid.UUID)

	var body CreateNoteRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if body.VideoID == "" || body.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "videoId and content required"})
		return
	}

	fmt.Printf("CAnt even fetch\n\n\n\n\n")

	note := database.Note{
		UserID:  userID,
		VideoID: body.VideoID,
		Content: body.Content,
		Tags:    body.Tags,
	}

	if err := database.InsertNote(&note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}

func GetNotes(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	videoId := c.Query("videoId")
	if videoId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "videoId required"})
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		fmt.Sscan(l, &limit)
	}

	var cursor *time.Time
	if cur := c.Query("cursor"); cur != "" {
		t, err := time.Parse(time.RFC3339, cur)
		if err == nil {
			cursor = &t
		}
	}

	tags := c.QueryArray("tag")

	var (
		notes []database.Note
		err   error
	)

	if len(tags) > 0 {
		notes, err = database.GetNotesByTags(videoId, limit, tags, cursor)
	} else {
		notes, err = database.GetNotes(videoId, limit, cursor)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var nextCursor *time.Time
	if len(notes) > 0 {
		t := notes[len(notes)-1].CreatedAt
		nextCursor = &t
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      notes,
		"nextCursor": nextCursor,
	})
}

func DeleteNote(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}
	userID := userIDAny.(uuid.UUID)

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&database.Note{})

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
