package router

import (
	"fmt"
	"go-todolist/controller"
	"go-todolist/entity"
	"go-todolist/middleware"
	"go-todolist/services"
	s3_utils "go-todolist/utils/aws"
	gorm_utils "go-todolist/utils/gorm"
	"go-todolist/utils/log"
	redis_utils "go-todolist/utils/redis"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	v1 = "/api/v1"

	db                    *gorm.DB                         = gorm_utils.InitMySQL()
	rdb                   *redis.Client                    = redis_utils.InitRedis()
	awsS3                 *s3.Client                       = s3_utils.InitS3()
	userEntity            entity.UserEntity                = entity.NewUserEntity(db)
	categoryEntity        entity.CategoryEntity            = entity.NewCategoryEntity(db)
	taskEntity            entity.TaskEntity                = entity.NewTaskEntity(db)
	redisEntity           entity.RedisEntity               = entity.NewRedisEntity(rdb)
	s3Entity              entity.S3Entity                  = entity.NewS3Entity(awsS3)
	userService           services.UserService             = services.NewUserService(userEntity)
	categoryService       services.CategoryService         = services.NewCategoryService(categoryEntity)
	taskService           services.TaskService             = services.NewTaskService(taskEntity, s3Entity)
	jwtService            services.JWTService              = services.NewJWTService(redisEntity, userEntity)
	userController                                         = controller.NewUserController(userService, jwtService)
	categoryController                                     = controller.NewCategoryController(categoryService, categoryEntity)
	taskController                                         = controller.NewTaskController(taskService, taskEntity)
	googleOauthController                                  = controller.NewGoogleOauthController(jwtService)
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
	// r.Use(middleware.CORS())

	// IPv6 0:0:0:0:0:0:0:1 = ::1 (Omit 0) = 0.0.0.0/0
	r.SetTrustedProxies([]string{"::1", "192.168.0.0/16", "172.16.0.0/12", "127.0.0.1/8", "10.0.0.0/8", "0.0.0.0/0"})

	// Set the IP rate limiter (limiter times, time)
	r.Use(rateLimiterMiddleware.RateLimiter(100, 60))

	authRoutes := r.Group(v1 + "/auth")
	{
		authRoutes.POST("/login", userController.Login)
		authRoutes.POST("/register", userController.Register)
	}

	oauthRoutes := r.Group(v1 + "/oauth")
	{
		oauthRoutes.GET("/google/login", controller.GoogleLogin)
		oauthRoutes.GET("/google/callback", googleOauthController.GoogleCallBack)
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

	categories := r.Group(v1+"/category", middleware.AuthorizeJWT(jwtService))
	{
		categories.POST("/", categoryController.Create)
		categories.GET("/", categoryController.GetByList)
		categories.GET("/:id", categoryController.Get)
		categories.PATCH("/:id", categoryController.Update)
		categories.DELETE("/:id", categoryController.Delete)
	}

	tasks := r.Group(v1+"/task", middleware.AuthorizeJWT(jwtService))
	{
		tasks.POST("/", taskController.Create)
		tasks.GET("/", taskController.GetByList)
		tasks.GET("/:id", taskController.Get)
		tasks.PATCH("/:id", taskController.Update)
		tasks.DELETE("/:id", taskController.Delete)
	}

	err := r.Run(appPort)
	if err != nil {
		return nil
	}

	return r
}
