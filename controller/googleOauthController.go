package controller

import (
	"fmt"
	"go-todolist/model"
	"go-todolist/services"
	"go-todolist/utils/responses"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOauthController interface {
	GoogleCallBack(c *gin.Context)
}

type googleOauthController struct {
	jwtService services.JWTService
}

func NewGoogleOauthController(jwtService services.JWTService) GoogleOauthController {
	return &googleOauthController{
		jwtService: jwtService,
	}
}

type GoogleInfo struct {
	Id    uint64
	Email string
}

var (
	googleOAuthConfig = &oauth2.Config{
		Endpoint:    google.Endpoint,
		RedirectURL: "http://localhost:8642/api/v1/oauth/google/callback",
		Scopes:      []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
	googleState = os.Getenv("GOOGLE_OAUTH_SECRET")
)

func GoogleLogin(c *gin.Context) {
	googleOAuthConfig.ClientID = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")

	url := googleOAuthConfig.AuthCodeURL(googleState)
	fmt.Println(url)
	c.Redirect(http.StatusSeeOther, url)
}

func (h *googleOauthController) GoogleCallBack(c *gin.Context) {
	var responseToken model.Token
	googleOAuthConfig.ClientSecret = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	code := c.Query("code")
	state := c.Query("state")
	if state != googleState {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to GoogleCallBack", responses.FailedToGetStateToken, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	token, err := googleOAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to GoogleCallBack", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to GoogleCallBack", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to GoogleCallBack", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	data := GoogleInfo{
		Id:    gjson.GetBytes(content, "id").Uint(),
		Email: gjson.GetBytes(content, "email").String(),
	}

	generatorToken := h.jwtService.GoogleGenerateToken(data)
	if len(generatorToken) < 1 {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to GoogleCallBack", responses.EmailNotExists, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	responseToken.Token = generatorToken

	response := responses.SuccessResponse(http.StatusOK, "Google access success", generatorToken)
	c.AbortWithStatusJSON(http.StatusOK, response)
	return
}
