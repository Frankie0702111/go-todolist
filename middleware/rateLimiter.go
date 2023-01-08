package middleware

import (
	"go-todolist/entity"
	"go-todolist/utils/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiterMiddleware interface {
	RateLimiter(limit int, slidingWindow int64) gin.HandlerFunc
}

type rateLimiterMiddleware struct {
	// conntection to redis
	redisEntity entity.RedisEntity
}

func NewRateLimiterMiddleware(redisEntity entity.RedisEntity) RateLimiterMiddleware {
	return &rateLimiterMiddleware{
		redisEntity: redisEntity,
	}
}

func (r *rateLimiterMiddleware) RateLimiter(limit int, expireAt int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now().Unix()

		getLimiter, _ := r.redisEntity.GetInt(ip)
		if getLimiter == limit {
			response := responses.ErrorsResponseByCode(http.StatusTooManyRequests, "Failed to get IP rate limiter", responses.TooManyRequests, nil)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response)
			return
		}

		r.redisEntity.IncrBy(ip, 1)
		if getLimiter == 0 {
			r.redisEntity.ExpireAt(ip, time.Unix(now+expireAt, 0))
		}

		c.Next()
	}
}
