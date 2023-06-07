package middleware

import (
	"go-todolist/controller"
	"go-todolist/services"
	"go-todolist/utils/responses"
	"log"
	"net/http"

	// logger "go-todolist/utils/log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(s services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := controller.ParseToken(c)
		switch authHeader {
		case "No token found":
			response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.NoTokenFound, nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		case "Bearer token not in proper format":
			response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.BearerTokenNotInProperFormat, nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Validate the token
		token, err := s.ValidateToken(authHeader)
		if err != nil {
			response := responses.ErrorsResponse(http.StatusUnauthorized, "Token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		} else {
			// Get the claims of the token
			claims := token.Claims.(jwt.MapClaims)
			// output the user_id
			log.Println("Claim[user_id]: ", claims["user_id"])
			// output the issuer
			log.Println("Claim[issuer] :", claims["iss"])
		}

		// whitelist for token
		redisToken := s.AuthJWT(authHeader)
		if len(redisToken) < 1 {
			response := responses.ErrorsResponseByCode(http.StatusUnauthorized, "Token is not valid", responses.TokenDoesNotExistOrExpired, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Check if the token is the same as redis
		if authHeader != redisToken {
			response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.TokenInvalid, nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		c.Next()
	}
}
