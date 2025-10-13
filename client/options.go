package client

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
)

// Options 配置
type Options struct {
	// 基础拨号
	Target    string   // 形如 "dns:///host:port" 或 "ipv4:port"；本地可直接 "127.0.0.1:9090"
	Endpoints []string // 可选：["10.0.0.1:9090","10.0.0.2:9090"]，内部自动 round_robin
	Authority string   // 可选：HTTP/2 :authority

	// 安全
	Insecure bool        // true 则明文；否则走 TLS
	TLSCfg   *tls.Config // 可选：自定义 TLS（可含 ServerName、RootCAs、mTLS 等）

	// 调用默认值
	DefaultTimeout time.Duration // 默认每次调用超时（建议 2~5s）

	// 拦截器
	UnaryInts  []grpc.UnaryClientInterceptor
	StreamInts []grpc.StreamClientInterceptor

	// 低级定制
	ExtraDialOptions []grpc.DialOption
	// gRPC ServiceConfig JSON：可用于开启官方重试/负载均衡策略等
	// 如：`{"loadBalancingPolicy":"round_robin","methodConfig":[{...retryPolicy...}]}`
	ServiceConfigJSON string
}

func (o *Options) fillDefaults() {
	if o.DefaultTimeout <= 0 {
		o.DefaultTimeout = 5 * time.Second
	}
}
