package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LXJ0000/todolist-grpc-gateway/config/bootstrap"
	"github.com/LXJ0000/todolist-grpc-gateway/discovery"
	service "github.com/LXJ0000/todolist-grpc-gateway/internal/service/pb"
	"github.com/LXJ0000/todolist-grpc-gateway/routes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	fmt.Println(env)

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	var slogHandler slog.Handler = slog.NewTextHandler(os.Stdout, opts)

	etcdAddress := []string{env.EtcdAddress}
	etcdRegister := discovery.NewResolver(etcdAddress, slog.New(slogHandler))
	resolver.Register(etcdRegister)
	go startListen()
	{
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		s := <-osSignal
		fmt.Println("exit", s)
	}
	fmt.Println("gateway listen on :4000")
}

func startListen() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	userConn, err := grpc.Dial("localhost:10001", opts...)
	if err != nil {
		panic(err)
	}
	userService := service.NewUserServiceClient(userConn)

	r := routes.NewRouter(userService)

	server := &http.Server{
		Addr:           ":4000",
		Handler:        r,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server start failed", "error", err)
	}
}
