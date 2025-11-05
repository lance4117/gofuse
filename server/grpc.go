package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// GrpcServer 封装一个可控生命周期的 gRPC 服务
type GrpcServer struct {
	opts   Options
	srv    *grpc.Server
	ln     net.Listener
	health *health.Server
}

// NewGrpc 创建 gRPC GrpcServer（不启动监听）
func NewGrpc(opts Options) (*GrpcServer, error) {
	if err := opts.fillDefaults(); err != nil {
		return nil, err
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(opts.unaryInterceptors()...),
		grpc.ChainStreamInterceptor(opts.streamInterceptors()...),
	}

	// TLS
	if opts.TLSCertFile != "" && opts.TLSKeyFile != "" {
		creds, err := loadServerTLS(opts.TLSCertFile, opts.TLSKeyFile, opts.TLSClientCAFile)
		if err != nil {
			return nil, fmt.Errorf("load tls: %w", err)
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	srv := grpc.NewServer(grpcOpts...)
	h := health.NewServer()

	// 注册健康检查
	healthpb.RegisterHealthServer(srv, h)

	// 注册业务服务
	if opts.Register != nil {
		if err := opts.Register(srv); err != nil {
			return nil, fmt.Errorf("register service: %w", err)
		}
	}

	// 开发环境可开启反射（方便 grpcui / grpcurl 调试）
	if opts.EnableReflection {
		reflection.Register(srv)
	}

	return &GrpcServer{
		opts:   opts,
		srv:    srv,
		health: h,
	}, nil
}

// Start 监听并启动服务；阻塞直到 ctx 取消或监听出错。
// 一般由外部用 goroutine 并发启动多个服务（HTTP/GRPC），统一等 ctx.Done().
func (s *GrpcServer) Start(ctx context.Context) error {
	var err error
	s.ln, err = net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.opts.Addr, err)
	}

	// 标记 Serving
	s.health.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	errCh := make(chan error, 1)
	go func() {
		errCh := errCh // shadow for clarity
		if serveErr := s.srv.Serve(s.ln); serveErr != nil && !errors.Is(serveErr, grpc.ErrServerStopped) {
			errCh <- serveErr
		}
	}()

	select {
	case <-ctx.Done():
		return s.shutdown(ctx.Err())
	case e := <-errCh:
		return e
	}
}

// Stop 立即关闭（不优雅）。一般只在优雅退出超时兜底时调用。
func (s *GrpcServer) Stop() {
	if s.srv != nil {
		s.srv.Stop()
	}
	if s.ln != nil {
		_ = s.ln.Close()
	}
}

// GracefulStop 在给定超时时间内优雅退出，否则强制 Stop。
func (s *GrpcServer) GracefulStop(timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		if s.srv != nil {
			s.srv.GracefulStop()
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		s.Stop()
	}
}

func (s *GrpcServer) shutdown(reason error) error {
	// 标记 NotServing
	s.health.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	// 优雅退出（默认 5s）
	to := s.opts.GracefulTimeout
	if to <= 0 {
		to = 5 * time.Second
	}
	s.GracefulStop(to)
	return reason
}

func loadServerTLS(certFile, keyFile, clientCA string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// 可选：双向 TLS（校验客户端）
	if clientCA != "" {
		caPEM, err := os.ReadFile(clientCA)
		if err != nil {
			return nil, err
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(caPEM) {
			return nil, io.ErrUnexpectedEOF
		}
		cfg.ClientCAs = cp
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return credentials.NewTLS(cfg), nil
}
