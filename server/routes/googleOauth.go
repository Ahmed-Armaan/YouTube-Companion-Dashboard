package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"yt_dashboard.com/database"
	"yt_dashboard.com/utils"
)

type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func GetCredentials(c *gin.Context) {
	code := c.Query("code")
	oauthErr := c.Query("error")

	if oauthErr != "" {
		c.JSON(500, gin.H{"error": oauthErr})
		return
	}

	if code == "" {
		c.JSON(400, gin.H{"error": "Missing code"})
		return
	}

	tokenResponse, err := getTokens(code)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("Acccess Token: %s\n", tokenResponse.AccessToken)

	userInfo, err := getUserInfo(tokenResponse.AccessToken)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("%v\n", userInfo)

	// inserting access tokoen into cache
	utils.InsertAccessToken(userInfo.Sub, tokenResponse.AccessToken)

	// inserting user and refresh token into Database
	id, err := database.InsertUser(userInfo.Sub, userInfo.Name, userInfo.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cant store user in DataBase",
		})
		return
	}

	err = database.InsertToken(id, tokenResponse.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cant store token in DataBase",
		})
		return
	}

	//	cookie := &http.Cookie{
	//		Name:     "session",
	//		Value:    userInfo.Sub,
	//		Path:     "/",
	//		HttpOnly: false,
	//		Secure:   false,
	//		SameSite: http.SameSiteLaxMode,
	//	}
	//
	// http.SetCookie(c.Writer, cookie)
	// c.Redirect(http.StatusFound, "http://localhost:5173/")
	claims := map[string]any{
		"sub": userInfo.Sub,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
	}
	jwtStr, err := utils.SignJwt(claims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not sign JWT",
		})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session",
		Value:    jwtStr,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	c.Redirect(http.StatusFound, os.Getenv("FRONTEND_URL")+"/channels")
}

func getTokens(code string) (*TokenResponse, error) {
	var tokenResponse TokenResponse
	reqUrl := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	data.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return &tokenResponse, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &tokenResponse, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return &tokenResponse, fmt.Errorf(
			"token exchange failed: %s | body: %s",
			res.Status,
			string(bodyBytes),
		)
	}

	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		return &tokenResponse, err
	}

	return &tokenResponse, nil
}

func getUserInfo(token string) (*GoogleUserInfo, error) {
	reqUrl := "https://www.googleapis.com/oauth2/v3/userinfo"
	var googleUserInfo GoogleUserInfo

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return &googleUserInfo, err
	}

	authorizationStr := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", authorizationStr)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &googleUserInfo, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &googleUserInfo, err
	}

	err = json.Unmarshal(body, &googleUserInfo)
	if err != nil {
		return &googleUserInfo, err
	}

	return &googleUserInfo, nil
}

func refreshToken(refreshToken string, userId string) (string, error) {
	var tokenResponse TokenResponse
	reqUrl := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf(
			"token exchange failed: %s | body: %s",
			res.Status,
			string(bodyBytes),
		)
	}

	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	accessToken := tokenResponse.AccessToken
	utils.InsertAccessToken(userId, accessToken)
	return accessToken, nil
}
