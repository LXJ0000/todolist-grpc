package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/LXJ0000/todolist-grpc/app/user/discovery"
	"github.com/LXJ0000/todolist-grpc/app/user/internal/handler"
	service "github.com/LXJ0000/todolist-grpc/app/user/internal/service/pb"
	"github.com/LXJ0000/todolist-grpc/app/user/config/bootstrap"
	"google.golang.org/grpc"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	fmt.Println(env)

	etcdAddress := []string{env.EtcdAddress}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	var slogHandler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	etcdRegister := discovery.NewRegister(etcdAddress, time.Second, slog.New(slogHandler))
	userNode := discovery.Server{
		Name: "user",
		Addr: "localhost:8081",
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定服务
	service.RegisterUserServiceServer(server, handler.NewUserService(app.Orm))

	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(userNode, time.Second*10); err != nil {
		panic(err)
	}
	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}
