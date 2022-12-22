package middleware

import (
	"go-todolist/services"
	"go-todolist/utils/responses"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

//AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the header of the request (if any)
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

		// Validate the token
		token, err := jwtService.ValidateToken(authHeader)

		if err != nil {
			response := responses.ErrorsResponse(401, "Token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return

		} else {
			// Get the claims of the token
			claims := token.Claims.(jwt.MapClaims)
			// output the user_id
			log.Println("Claim[user_id]: ", claims["user_id"])
			// output the issuer
			log.Println("Claim[issuer] :", claims["issuer"])
		}
	}
}
