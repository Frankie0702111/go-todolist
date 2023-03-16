package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ----- First version -----
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, UPDATE, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Accept-Encoding, Authorization, Connection, Content-Length, Content-Type, Host, User-Agent, Token, Access-Control-Request-Method, Access-Control-Request-Headers")
		// c.Header("Access-Control-Allow-Credentials", "true") // cookie - If you need it, open it
		c.Header("Access-Control-Max-Age", "9000")
		c.Set("content-type", "application/json")

		// ----- Second version -----
		// origin := c.Request.Header.Get("Origin")
		// if origin != "" {
		// 	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		// 	c.Header("Access-Control-Allow-Origin", "*")
		// 	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, UPDATE, DELETE")
		// 	c.Header("Access-Control-Allow-Headers", "Accept, Accept-Encoding, Authorization, Connection, Content-Length, Content-Type, Host, User-Agent, Token, Access-Control-Request-Method, Access-Control-Request-Headers")
		// 	c.Header("Access-Control-Allow-Credentials", "true")
		// 	c.Header("Access-Control-Max-Age", "9000")
		// 	c.Set("content-type", "application/json")
		// }

		// ----- Old version -----
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, UPDATE, DELETE")
		// c.Writer.Header().Set("Access-Control-Allow-Max-Age", "9000")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
