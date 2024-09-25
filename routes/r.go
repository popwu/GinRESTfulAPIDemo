package routes

import (
	"github.com/gin-gonic/gin"

	api "apidemo/handlers/api"
	auth "apidemo/handlers/auth"
	middleware "apidemo/middleware"
)

func SetupRouter(r *gin.Engine) {
	r_auth := r.Group("/auth")
	{
		r_auth.POST("/login", auth.Login)
	}
	r_api := r.Group("/api")
	{
		r_api.Use(middleware.JwtAuthMiddleware())
		r_api.GET("/users", api.GetUsers)
	}

}
