package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/iam-jessicawong/mygram/controllers"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/user")
	{
		userRouter.POST("/register", controllers.Register)
		userRouter.POST("/login", controllers.Login)
	}

	return r
}
