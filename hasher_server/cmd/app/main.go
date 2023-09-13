package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	formatters "github.com/fabienm/go-logrus-formatters"
	"github.com/sirupsen/logrus"

	"hasher_server/internal/handler"
	"proto/hash_service"

	"google.golang.org/grpc"
)

func main() {
	var gelFmt = formatters.NewGelf("client")
	var Logger = logrus.Logger{
		Out:       os.Stdout,
		Formatter: gelFmt,
		Level:     logrus.InfoLevel,
		Hooks:     make(logrus.LevelHooks),
	}
	Logger.Info("Starting server ...")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	port := os.Getenv("HASH_SERVICE_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		Logger.Fatal("failed to listen:", err)
	}
	s := grpc.NewServer()

	server := &handler.Server{Logger: &Logger}
	hash_service.RegisterHashServiceServer(s, server)
	go func() {
		if err := s.Serve(lis); err != nil {
			Logger.Fatal("failed to serve:", err)
		}
	}()
	<-sigCh
	s.GracefulStop()
	Logger.Info("Server stopped")
	os.Exit(0)
}
