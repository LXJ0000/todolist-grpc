package handler

import (
	"context"

	service "github.com/LXJ0000/todolist-grpc-gateway/internal/service/pb"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var userReq service.UserRequest
	if err := c.ShouldBind(&userReq); err != nil {
		// TODO: 处理错误
		return
	}
	in := c.Keys["user"]
	userService := in.(service.UserServiceClient)
	userResp, err := userService.Register(context.Background(), &userReq)
	if err != nil {
		// TODO: 处理错误
		return
	}
	c.JSON(200, userResp) // TODO: 完善
}

func Login(c *gin.Context) {
	var userReq service.UserRequest
	if err := c.ShouldBind(&userReq); err != nil {
		// TODO: 处理错误
		return
	}
	// 获取服务实例
	userService := c.Keys["user"].(service.UserServiceClient)
	userResp, err := userService.Login(context.Background(), &userReq)
	// TODO: JWT
	if err != nil {
		// TODO: 处理错误
		return
	}
	c.JSON(200, userResp) // TODO: 完善
}
