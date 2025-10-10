package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lance4117/gofuse/logger"
	"github.com/lance4117/gofuse/server"
	"github.com/lance4117/gofuse/server/test/pb"
	"google.golang.org/grpc"
)

func TestInitServer(t *testing.T) {
	engine := server.NewHTTP(true)
	engine.POST("/hello", Query1())

	engine.Run()
}

func Query1() server.ContextHandler {
	return func(ctx *server.Context) {
		ctx.Response(200, "hello world")
	}
}

func TestGrpc(t *testing.T) {
	gs, err := server.NewGrpc(server.Options{Addr: ":9090",
		Register: register,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = gs.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func register(s *grpc.Server) error {
	pb.RegisterDemoServer(s, &DemoService{})

	return nil
}

// DemoService 实现 pb.DemoServer 接口
type DemoService struct {
	pb.UnimplementedDemoServer // 必须嵌入，避免接口变动报错
}

// Unary：普通请求
func (s *DemoService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	logger.Infof("[SayHello] name=%s meta=%v", req.Name, req.Meta)
	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
		TsUnix:  time.Now().Unix(),
	}, nil
}
