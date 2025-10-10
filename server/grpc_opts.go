package server

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/lance4117/gofuse/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Options 服务配置
type Options struct {
	Addr             string
	TLSCertFile      string // 可选：服务端证书
	TLSKeyFile       string // 可选：服务端私钥
	TLSClientCAFile  string // 可选：开启双向 TLS 时的 Client CA
	EnableReflection bool
	GracefulTimeout  time.Duration
	UnaryInts        []grpc.UnaryServerInterceptor
	StreamInts       []grpc.StreamServerInterceptor

	// Register 用于注册你的具体 gRPC 服务（pb.RegisterXxxServer）
	Register func(s *grpc.Server) error
}

func (o *Options) fillDefaults() error {
	if o.Addr == "" {
		o.Addr = ":9090"
	}
	return nil
}

func (o *Options) unaryInterceptors() []grpc.UnaryServerInterceptor {
	ints := []grpc.UnaryServerInterceptor{unaryRecovery(), unaryLogging()}
	ints = append(ints, o.UnaryInts...)
	return ints
}

func (o *Options) streamInterceptors() []grpc.StreamServerInterceptor {
	ints := []grpc.StreamServerInterceptor{streamRecovery(), streamLogging()}
	ints = append(ints, o.StreamInts...)
	return ints
}

// 简单日志拦截器
func unaryLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		st, _ := status.FromError(err)
		logger.Infof("[grpc] unary %s code=%v cost=%s", info.FullMethod, st.Code(), time.Since(start))
		return resp, err
	}
}

func streamLogging() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		st, _ := status.FromError(err)
		logger.Infof("[grpc] stream %s code=%v cost=%s", info.FullMethod, st.Code(), time.Since(start))
		return err
	}
}

// panic 恢复拦截器
func unaryRecovery() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Infof("[grpc] panic recovered: %v\n%s", r, debug.Stack())
				err = status.Error(13, "internal") // codes.Internal
			}
		}()
		return handler(ctx, req)
	}
}

func streamRecovery() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Infof("[grpc] panic recovered: %v\n%s", r, debug.Stack())
				err = status.Error(13, "internal")
			}
		}()
		return handler(srv, ss)
	}
}
