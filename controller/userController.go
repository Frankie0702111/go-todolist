package controller

import (
	"go-todolist/model"
	"go-todolist/request"
	"go-todolist/services"
	"go-todolist/utils/responses"
	"net/http"
	"strconv"
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

// Login is a function for login
func (h *userController) Login(c *gin.Context) {
	// create new instance of LoginRequest
	var input request.LoginRequest
	jwtTTL := services.GetTokenTTL()

	// bind the input with the request body
	err := c.ShouldBindJSON(&input)

	// Check if there is any error in binding
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), responses.EmptyObject{})
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Check if the email and password is valid
	loginResult := h.userService.VerifyCredential(input.Email, input.Password)
	if v, ok := loginResult.(model.User); ok {
		generatedToken := h.jwtService.GenerateToken(strconv.Itoa(int(v.ID)), time.Now().Add(time.Duration(jwtTTL)*time.Second))

		if len(generatedToken) < 1 {
			response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", "Signature failed", responses.EmptyObject{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		v.Token = generatedToken
		response := responses.SuccessResponse(http.StatusOK, "Login successfully", v)
		c.JSON(http.StatusOK, response)
		return
	}

	// If the email and password is not valid
	response := responses.ErrorsResponse(401, "Failed to process request", "Invalid credential", responses.EmptyObject{})
	c.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

// Register is a function for register
func (h *userController) Register(c *gin.Context) {
	// create new instance of RegisterRequest
	var input request.RegisterRequest

	// bind the register with the request body
	err := c.ShouldBind(&input)

	// Check if there is any error in binding
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), responses.EmptyObject{})
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// if the email is valid and unique in the database then register the user
	// create new user
	createdUser := h.userService.CreateUser(input)

	if len(createdUser.Password) < 1 {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", "Failed to hash a password", responses.EmptyObject{})
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// response with the user data and token
	response := responses.SuccessResponse(http.StatusCreated, "Register Success", createdUser)

	// return the response
	c.JSON(http.StatusCreated, response)
}

func (h *userController) RefreshToken(c *gin.Context) {
	var refreshToken model.RefreshToken

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		response := responses.ErrorsResponse(401, "Failed to process request", "No token found", nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		response := responses.ErrorsResponse(401, "Failed to process request", "Bearer token not in proper format", nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authHeader = strings.TrimSpace(splitToken[1])

	// Refresh the token
	token := h.jwtService.RefreshToken(authHeader)

	if len(token) < 1 {
		response := responses.ErrorsResponse(401, "Failed to process request", "Token contains an invalid number of segments", nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	refreshToken.Token = token
	response := responses.SuccessResponse(http.StatusOK, "Refresh token successfully", refreshToken)
	c.JSON(http.StatusOK, response)
	return
}

func (h *userController) Logout(c *gin.Context) {
	// var refreshToken model.RefreshToken
	// authHeader := c.GetHeader("Authorization")

	// if authHeader == "" {
	// 	response := responses.ErrorsResponse(401, "Failed to process request", "No token found", nil)
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, response)
	// 	return
	// }

	// splitToken := strings.Split(authHeader, "Bearer ")
	// if len(splitToken) != 2 {
	// 	response := responses.ErrorsResponse(401, "Failed to process request", "Bearer token not in proper format", nil)
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, response)
	// 	return
	// }
	// authHeader = strings.TrimSpace(splitToken[1])

	// // Refresh the token
	// token := h.jwtService.GenerateToken(authHeader, time.Now())

	// refreshToken.Token = token
	response := responses.SuccessResponse(http.StatusOK, "Successfully logged out", nil)
	c.JSON(http.StatusOK, response)
	return
}
