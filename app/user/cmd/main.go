package main

import (
	"fmt"

	"github.com/LXJ0000/todolist-grpc/config/bootstrap"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	fmt.Println(env)

	// db := app.Orm
	// cache := app.Cache

	// timeout := time.Duration(env.ContextTimeout) * time.Minute // 接口超时时间

	// server := gin.Default()

	// _ = server.Run(env.ServerAddress)
}
