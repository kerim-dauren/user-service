package main

import (
	"context"
	"errors"
	"fmt"
	user "github.com/kerim-dauren/user-service/gen/proto"
	"github.com/kerim-dauren/user-service/internal/api"
	v1 "github.com/kerim-dauren/user-service/internal/api/grpc/v1"
	"github.com/kerim-dauren/user-service/internal/configs"
	"github.com/kerim-dauren/user-service/internal/services"
	"github.com/kerim-dauren/user-service/internal/storages/pg"
	"github.com/kerim-dauren/user-service/pkg/hashx"
	"github.com/kerim-dauren/user-service/pkg/postgresx"
	"github.com/kerim-dauren/user-service/pkg/slogx"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title User Service
// @version 1.0
// @description User Service API
// @BasePath /api/v1
// @schemes http https
func main() {
	ctx := context.Background()
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load app configs", err)
	}

	// FYI: A logger interface should ideally be created and implemented.
	// For now, using a pre-configured slog instance as a placeholder.
	logger, err := slogx.NewLogger(&slogx.Config{
		Level:   cfg.Log.Level,
		Handler: cfg.Log.Handler,
		Writer:  cfg.Log.Writer,
	})

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	dbPool, err := postgresx.New(cfg.DbUrl)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer dbPool.Close()

	userStorage := pg.NewUserStorage(dbPool)
	hasher := hashx.NewArgon2Hasher()
	userService := services.NewUserService(logger, userStorage, hasher)

	httpRouter := api.NewHttpRouter(&api.RouterDeps{
		UserService: userService,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: httpRouter,
	}

	errch := make(chan error, 1)

	// Starting the HTTP server
	go func() {
		logger.Info("http server started", "port", cfg.HttpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errch <- fmt.Errorf("failed to start http server: %w", err)
		}
		close(errch)
	}()

	// Starting the gRPC server
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
		if err != nil {
			errch <- fmt.Errorf("failed to listen grpc port: %w", err)
		}
		logger.Info("grpc server started", "port", cfg.GRPCPort)

		grpcServer := grpc.NewServer()
		user.RegisterUserServiceServer(grpcServer, v1.NewUserService(userService))

		if err := grpcServer.Serve(lis); err != nil {
			errch <- fmt.Errorf("failed to serve grpc server: %w", err)
		}
	}()

	select {
	case err := <-errch:
		log.Fatalln(err)
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = server.Shutdown(timeout); err != nil {
			log.Fatalln(err)
		}
	}
}
