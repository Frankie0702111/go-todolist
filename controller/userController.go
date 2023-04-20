package controller

import (
	"go-todolist/model"
	"go-todolist/request"
	"go-todolist/services"
	"go-todolist/utils/responses"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// User Controller interface is a contract for all user controller
type UserController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

// User Controller struct to implement UserController interface
type userController struct {
	// inject user service
	userService services.UserService

	// inject jwt service
	jwtService services.JWTService
}

// Create a new instance of UserController with userService and jwtService injected as dependency
func NewUserController(userService services.UserService, jwtService services.JWTService) UserController {
	return &userController{
		// inject user service
		userService: userService,

		// inject jwt service
		jwtService: jwtService,
	}
}

// ParseToken is a shared method for parse user token
func ParseToken(c *gin.Context) string {
	// Get the token from the header of the request (if any)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "No token found"
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "Bearer token not in proper format"
	}

	authHeader = strings.TrimSpace(splitToken[1])

	return authHeader
}

// Login is a function for user login
// @Summary "User Login"
// @Tags	"Auth"
// @Version 1.0
// @Produce application/json
// @Param	* body request.LoginRequest true "User Login"
// @Success 200 object responses.Response{errors=string,data=string} "Login successfully"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 401 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/auth/login [post]
func (h *userController) Login(c *gin.Context) {
	// create new instance of LoginRequest
	var input request.LoginRequest
	jwtTTL := services.GetTokenTTL()

	// bind the input with the request body
	err := c.ShouldBindJSON(&input)
	// Check if there is any error in binding
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		// response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), responses.EmptyObject{})
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Check if the email and password is valid
	loginResult := h.userService.VerifyCredential(input.Email, input.Password)
	if v, ok := loginResult.(model.User); ok {
		generatedToken := h.jwtService.GenerateToken(v.ID, time.Now().Add(time.Duration(jwtTTL)*time.Second))
		if len(generatedToken) < 1 {
			response := responses.ErrorsResponseByCode(http.StatusInternalServerError, "Failed to process request", responses.SignatureFailed, nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		v.Token = generatedToken
		response := responses.SuccessResponse(http.StatusOK, "Login successfully", v)
		c.JSON(http.StatusOK, response)
		return
	}

	// If the email and password is not valid
	response := responses.ErrorsResponseByCode(http.StatusUnauthorized, "Failed to process request", responses.InvalidCredential, nil)
	c.AbortWithStatusJSON(http.StatusUnauthorized, response)
	return
}

// Register is a function for user register
// @Summary "User Register"
// @Tags	"Auth"
// @Version 1.0
// @Produce application/json
// @Param	* body request.RegisterRequest true "User Register"
// @Success 201 object responses.Response{errors=string,data=string} "Register Success"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/auth/register [post]
func (h *userController) Register(c *gin.Context) {
	// create new instance of RegisterRequest
	var input request.RegisterRequest

	// bind the register with the request body
	err := c.ShouldBind(&input)
	// Check if there is any error in binding
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// if the email is valid and unique in the database then register the user
	// create new user
	createdUser := h.userService.CreateUser(input)
	// Check if then email exists
	if createdUser.ID == 0 {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.EmailAlreadyExists, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// response with the user data and token
	response := responses.SuccessResponse(http.StatusCreated, "Register Success", createdUser)
	// return the response
	c.JSON(http.StatusCreated, response)
	return
}

// RefreshToken is a function for token refresh
// @Summary "User Refresh Token"
// @Tags	"Auth"
// @Version 1.0
// @Produce application/json
// @Param	Authorization header string true "example:Bearer token (Bearer+space+token)." default(Bearer )
// @Success 200 object responses.Response{errors=string,data=string} "Refresh token successfully"
// @Failure 401 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/auth/refresh [post]
func (h *userController) RefreshToken(c *gin.Context) {
	var refreshToken model.Token
	authHeader := ParseToken(c)
	token := h.jwtService.RefreshToken(authHeader)
	if len(token) < 1 {
		response := responses.ErrorsResponseByCode(http.StatusUnauthorized, "Failed to process request", responses.TokenContainsAnInvalidNumberOfSegments, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// set up return format
	refreshToken.Token = token
	response := responses.SuccessResponse(http.StatusOK, "Refresh token successfully", refreshToken)
	c.JSON(http.StatusOK, response)
	return
}

// Logout is a function for user logout
// @Summary "User Logout"
// @Tags	"Auth"
// @Version 1.0
// @Produce application/json
// @Param	Authorization header string true "example:Bearer token (Bearer+space+token)." default(Bearer )
// @Success 200 object responses.Response{errors=string,data=string} "Successfully logged out"
// @Failure 401 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/auth/logout [post]
func (h *userController) Logout(c *gin.Context) {
	authHeader := ParseToken(c)
	logout := h.jwtService.Logout(authHeader)
	if logout == true {
		response := responses.SuccessResponse(http.StatusOK, "Successfully logged out", nil)
		c.JSON(http.StatusOK, response)
	} else {
		response := responses.ErrorsResponseByCode(http.StatusUnauthorized, "Failed to process request", responses.FailedToLogout, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}

	return
}
