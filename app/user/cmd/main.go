package main

import (
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/LXJ0000/todolist-grpc-user/config/bootstrap"
	"github.com/LXJ0000/todolist-grpc-user/discovery"
	"github.com/LXJ0000/todolist-grpc-user/internal/handler"
	service "github.com/LXJ0000/todolist-grpc-user/internal/service/pb"
	"google.golang.org/grpc"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	etcdAddress := []string{env.EtcdAddress}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	var slogHandler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	etcdRegister := discovery.NewRegister(etcdAddress, time.Second*3, slog.New(slogHandler))
	userNode := discovery.Server{
		Name: "user",
		Addr: "localhost:10001",
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定服务
	service.RegisterUserServiceServer(server, handler.NewUserService(app.Orm))

	listen, err := net.Listen("tcp", "localhost:10001")
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
