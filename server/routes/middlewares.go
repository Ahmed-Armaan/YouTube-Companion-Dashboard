package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"yt_dashboard.com/database"
	"yt_dashboard.com/utils"
)

func VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// extracting the user Sub (userId) from the cookies send using JWT
		cookie, err := c.Cookie("session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "No authentication token provided",
			})
			return
		}

		claims, err := utils.VerifyJwt(cookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Cant authenticate provided token",
			})
			return
		}

		userId, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Token information unavailable",
			})
			return
		}

		// access token Resolution flow:
		// in cache?
		// ├─Yes -> use access token
		// └─No
		//    └─in DB? (refresh token)
		//      ├─Yes -> refresh access token
		//      └─No -> Error

		token, err := checkCache(userId)
		if err != nil {
			refreshToken_, err := database.GetToken(userId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}

			token, err = refreshToken(refreshToken_, userId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		c.Set("accessToken", token)
		c.Next()
	}
}

func checkCache(userID string) (string, error) {
	token, exists := utils.GetAccessTokenFromCache(userID)
	if !exists {
		return "", errors.New("No access token availbale in cache")
	} else {
		return token, nil
	}
}
