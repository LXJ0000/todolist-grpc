package routes

import (
	"github.com/LXJ0000/todolist-grpc-gateway/internal/handler"
	"github.com/LXJ0000/todolist-grpc-gateway/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(service ...interface{}) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Init(service)) // TODO: CROS JWT
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "pong",
			})
		})
		// User
		v1.POST("/user/login", handler.Login)
		v1.POST("/user/register", handler.Register)
		// Task
	}
	return r
}
