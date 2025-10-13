package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type GrpcClient[T any] struct {
	conn *grpc.ClientConn
	Stub T // 具体业务客户端（由 ctor 构造）
}

// NewGrpcClient 通过 ctor 构造具体的 pb 客户端（例如 pb.NewDemoClient）
func NewGrpcClient[T any](ctor func(cc *grpc.ClientConn) T, opts Options) (*GrpcClient[T], error) {
	opts.fillDefaults()

	target := opts.Target
	dialOpts := make([]grpc.DialOption, 0, 8)

	// 负载均衡：若提供多个 Endpoints，则采用 round_robin
	if len(opts.Endpoints) > 0 && target == "" {
		// 使用“解析器 + round_robin”是正道，这里给一个快速可用的统一 target：
		// 通过 passthrough + 多地址不太通用；更稳的做法是自行注册自定义 resolver。
		// 为简洁，这里推荐直接用第一个地址；或要求你用 DNS 做服务发现。
		// 若你确实需要多地址，下面开启 round_robin 并要求通过 DNS 返回多 A 记录。
		target = opts.Endpoints[0]
		dialOpts = append(dialOpts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"`+roundrobin.Name+`"}`))
	}

	if target == "" {
		return nil, fmt.Errorf("grpcx: empty Target/Endpoints")
	}

	// 凭证
	if opts.Insecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		var cfg *tls.Config
		if opts.TLSCfg != nil {
			cfg = opts.TLSCfg.Clone()
		} else {
			cfg = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(cfg)))
	}

	// ServiceConfig（重试/超时/负载均衡等）
	if sc := strings.TrimSpace(opts.ServiceConfigJSON); sc != "" {
		dialOpts = append(dialOpts, grpc.WithDefaultServiceConfig(sc))
	}

	// 拦截器
	if len(opts.UnaryInts) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(opts.UnaryInts...))
	}
	if len(opts.StreamInts) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainStreamInterceptor(opts.StreamInts...))
	}

	// 其他
	if opts.Authority != "" {
		dialOpts = append(dialOpts, grpc.WithAuthority(opts.Authority))
	}
	dialOpts = append(dialOpts, opts.ExtraDialOptions...)

	cc, err := grpc.NewClient(target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("grpc dial %s: %w", target, err)
	}

	return &GrpcClient[T]{conn: cc, Stub: ctor(cc)}, nil
}

func (c *GrpcClient[T]) Conn() *grpc.ClientConn { return c.conn }

func (c *GrpcClient[T]) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Do 包一层默认超时，统一上下文
func (c *GrpcClient[T]) Do(ctx context.Context, defaultTimeout time.Duration, fn func(ctx context.Context, stub T) error) error {
	if defaultTimeout <= 0 {
		defaultTimeout = 5 * time.Second
	}
	var cancel context.CancelFunc
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
		defer cancel()
	}
	return fn(ctx, c.Stub)
}

// WaitReady 使用 gRPC Health Check 等待服务就绪
func (c *GrpcClient[T]) WaitReady(ctx context.Context, svc string) error {
	hc := healthpb.NewHealthClient(c.conn)
	_, err := hc.Check(ctx, &healthpb.HealthCheckRequest{Service: svc})
	return err
}
