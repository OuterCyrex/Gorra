package Gorra

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func ServerRun(server *grpc.Server, serverConfig ServerConfig) error {
	port, err := getFreePort()
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverConfig.Addr, port))
	if err != nil {
		return err
	}

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	err = HealthCheck(serverConfig, fmt.Sprintf("%s:%d", serverConfig.Addr, port), 15)
	if err != nil {
		return err
	}

	fmt.Printf("\033[31m[Gorra] Server Runs On Port %d......\033[0m\n", port)

	err = server.Serve(lis)
	if err != nil {
		return err
	}
	defer server.GracefulStop()

	return nil
}
