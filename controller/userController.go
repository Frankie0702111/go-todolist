package controller

import (
	"go-todolist/model"
	"go-todolist/request"
	"go-todolist/services"
	"go-todolist/utils/responses"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// User Controller interface is a contract for all user controller
type UserController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
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
		generatedToken := h.jwtService.GenerateToken(strconv.Itoa(int(v.ID)))

		if len(generatedToken) < 1 {
			response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", "Signature failed", responses.EmptyObject{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		v.Token = generatedToken
		response := responses.SuccessResponse(http.StatusOK, "Login Successfully", v)
		c.JSON(http.StatusOK, response)
		return
	}

	// If the email and password is not valid
	response := responses.ErrorsResponse(401, "Failed to process request", "Invalid Credential", responses.EmptyObject{})
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
