package test

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func TestGrpc(t *testing.T) {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	stream, err := client.ServerReflectionInfo(context.Background())
	if err != nil {
		panic(err)
	}

	// 请求服务列表
	err = stream.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{},
	})
	if err != nil {
		panic(err)
	}

	resp, err := stream.Recv()
	if err != nil {
		panic(err)
	}

	for _, svc := range resp.GetListServicesResponse().Service {
		fmt.Println("Service:", svc.Name)
	}

}
