package grpc

import (
	"context"
	"fmt"
	"net"

	"backend/pkg/config"
	"backend/shared/logger"

	"google.golang.org/grpc"
	//"backend/internal/shared/trace"
)

type Server struct {
	config     *config.Config
	logger     logger.Logger
	grpcServer *grpc.Server
}

func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
	// grpc.UnaryInterceptor(trace.UnaryServerInterceptor()),
	// Add more interceptors here
	)

	return &Server{
		config:     cfg,
		logger:     logger,
		grpcServer: grpcServer,
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.config.GRPC.Host, s.config.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info(ctx, "Starting gRPC server on "+addr, map[string]interface{}{})
	return s.grpcServer.Serve(listener)
}

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}
