package router

import (
	"fmt"
	"go-todolist/controller"
	"go-todolist/entity"
	"go-todolist/middleware"
	"go-todolist/services"
	gorm_utils "go-todolist/utils/gorm"
	"go-todolist/utils/log"
	redis_utils "go-todolist/utils/redis"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	v1 = "/api/v1"

	db                    *gorm.DB                         = gorm_utils.InitMySQL()
	rdb                   *redis.Client                    = redis_utils.InitRedis()
	userEntity            entity.UserEntity                = entity.NewUserEntity(db)
	redisEntity           entity.RedisEntity               = entity.NewRedisEntity(rdb)
	userService           services.UserService             = services.NewUserService(userEntity)
	jwtService            services.JWTService              = services.NewJWTService(redisEntity)
	userController                                         = controller.NewUserController(userService, jwtService)
	rateLimiterMiddleware middleware.RateLimiterMiddleware = middleware.NewRateLimiterMiddleware(redisEntity)
)

func SetupRouter() *gin.Engine {
	// Load .env file
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Panic("Failed to load env file")
	}

	appPort := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	// Closing the database when the program stop
	defer gorm_utils.Close(db)
	defer redis_utils.Close(rdb)

	// r := gin.New()
	r := gin.Default()
	r.SetTrustedProxies(nil)
	// Set the IP rate limiter (limiter times, time)
	r.Use(rateLimiterMiddleware.RateLimiter(100, 60))

	authRoutes := r.Group(v1 + "/auth")
	{
		authRoutes.POST("/login", userController.Login)
		authRoutes.POST("/register", userController.Register)
	}

	test := r.Group(v1+"/test", middleware.AuthorizeJWT(jwtService))
	{
		test.GET("/token", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Token verification success.",
			})
		})
	}

	auth := r.Group(v1+"/auth", middleware.AuthorizeJWT(jwtService))
	{
		auth.POST("refresh", userController.RefreshToken)
		auth.POST("logout", userController.Logout)
	}

	err := r.Run(appPort)
	if err != nil {
		return nil
	}

	return r
}
