package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"yt_dashboard.com/database"
	"yt_dashboard.com/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No env file found") // no file in prod
	}

	database.DbInit()
	runServer()
}

func runServer() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/auth/callback", routes.GetCredentials)
	r.GET("/me", routes.VerifyUser(), routes.Me)
	r.GET("/channelId", routes.VerifyUser(), routes.MyChannelId)
	//r.GET("/logout", routes.Logout)
	r.GET("/channel", routes.VerifyUser(), routes.GetChannel)
	r.GET("/comments", routes.VerifyUser(), routes.GetCommentThread)
	//r.PUT("/video/description", routes.VerifyUser(), routes.UpdateVideoDescription)
	//r.PUT("/video/title", routes.VerifyUser(), routes.UpdateVideoTitle)
	r.POST("/comments", routes.VerifyUser(), routes.AddComment)
	r.POST("/comments/reply", routes.VerifyUser(), routes.ReplyToComment)
	r.DELETE("/comments", routes.VerifyUser(), routes.DeleteComment)
	r.POST("/ai/title", routes.VerifyUser(), routes.SuggestTitles)
	r.POST("/notes", routes.VerifyUser(), routes.CreateNote)
	r.GET("/notes", routes.VerifyUser(), routes.GetNotes)
	r.DELETE("/notes", routes.VerifyUser(), routes.DeleteNote)
	r.Run(":3000")
}
