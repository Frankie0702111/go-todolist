package router

import (
	"fmt"
	"go-todolist/controller"
	"go-todolist/entity"
	"go-todolist/middleware"
	"go-todolist/services"
	gorm_utils "go-todolist/utils/gorm"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	v1 = "/api/v1"

	db             *gorm.DB             = gorm_utils.InitMySQL()
	userEntity     entity.UserEntity    = entity.NewUserEntity(db)
	userService    services.UserService = services.NewUserService(userEntity)
	jwtService     services.JWTService  = services.NewJWTService()
	userController                      = controller.NewUserController(userService, jwtService)
)

func SetupRouter() *gin.Engine {
	// Load .env file
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	appPort := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

	// Closing the database when the program stop
	defer gorm_utils.Close(db)

	// r := gin.New()
	r := gin.Default()
	r.SetTrustedProxies(nil)

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
