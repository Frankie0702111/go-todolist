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

//AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(s services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := controller.ParseToken(c)

		// Validate the token
		token, err := s.ValidateToken(authHeader)
		if err != nil {
			response := responses.ErrorsResponse(401, "Token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		} else {
			// Get the claims of the token
			claims := token.Claims.(jwt.MapClaims)
			// output the user_id
			log.Println("Claim[user_id]: ", claims["user_id"])
			// output the issuer
			log.Println("Claim[issuer] :", claims["issuer"])
		}

		// whitelist for token
		redisToken := s.AuthJWT(authHeader)
		if len(redisToken) < 1 {
			response := responses.ErrorsResponse(401, "Token is not valid", "Token does not exist or expired", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}

		return
	}
}
