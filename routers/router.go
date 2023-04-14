package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/iam-jessicawong/mygram/controllers"
	"github.com/iam-jessicawong/mygram/middlewares"

	_ "github.com/iam-jessicawong/mygram/docs"

	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

// @title MyGram API
// @version 1.0
// description This is a simple service for managing mygram
// @termOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/license/LICENSE-2.0.html
// @host https://mygram-production-83b7.up.railway.app
// @BasePath /
func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/user")
	{
		userRouter.POST("/register", controllers.Register)
		userRouter.POST("/login", controllers.Login)
	}

	photoRouter := r.Group("/photo")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.GET("/", controllers.GetAllPhoto)
		photoRouter.GET("/:id", controllers.GetOnePhoto)
		photoRouter.POST("/", controllers.CreatePhoto)
		photoRouter.PUT("/:id", middlewares.Authorization(), controllers.UpdatePhoto)
		photoRouter.DELETE("/:id", middlewares.Authorization(), controllers.DeletePhoto)
	}

	commentRouter := r.Group("/comment")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.GET("/", controllers.GetAllComment)
		commentRouter.GET("/:id", controllers.GetOneComment)
		commentRouter.POST("/:photoId", controllers.CreateComment)
		commentRouter.PUT("/:id", middlewares.Authorization(), controllers.UpdateComment)
		commentRouter.DELETE("/:id", middlewares.Authorization(), controllers.DeleteComment)
	}

	socialMediaRouter := r.Group("/social-media")
	{
		socialMediaRouter.Use(middlewares.Authentication())
		socialMediaRouter.GET("/", controllers.GetAllSocialMedia)
		socialMediaRouter.GET("/:id", controllers.GetOneSocialMedia)
		socialMediaRouter.POST("/", controllers.CreateSocialMedia)
		socialMediaRouter.PUT("/:id", middlewares.Authorization(), controllers.UpdateSocialMedia)
		socialMediaRouter.DELETE("/:id", middlewares.Authorization(), controllers.DeleteSocialMedia)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
