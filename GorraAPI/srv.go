package GorraAPI

import (
	"errors"
	"fmt"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetSrvConnection(wait uint, config BaseConfig, srvName string) (*grpc.ClientConn, error) {
	Conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s:%d/%s?wait=%ds",
			config.GetRegistryInfo().Consul.Host,
			config.GetRegistryInfo().Consul.Port,
			srvName,
			wait,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy": "%s"}`, "round_robin")),
		grpc.WithUnaryInterceptor(middlewares.GrpcTracerInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		msg := fmt.Sprintf("[Gorra]: Load Balance Init Failed: %v", err)
		return nil, errors.New(msg)
	}

	return Conn, nil
}
