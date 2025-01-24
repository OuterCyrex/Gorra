package Gorra

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func HealthCheck(serverConfig ServerConfig, grpcAddr string, checkInterval uint) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", serverConfig.Consul.Host, serverConfig.Consul.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	check := &api.AgentServiceCheck{
		GRPC:                           grpcAddr,
		Timeout:                        "5s",
		Interval:                       fmt.Sprintf("%ds", checkInterval),
		DeregisterCriticalServiceAfter: "15s",
	}

	addr := strings.Split(grpcAddr, ":")
	port, _ := strconv.Atoi(addr[1])

	serviceUUID := uuid.NewV4().String()

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    serverConfig.Name,
		ID:      serviceUUID,
		Port:    port,
		Tags:    serverConfig.Tags,
		Address: serverConfig.Addr,
		Check:   check,
	})

	if err != nil {
		return err
	}

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		if err = client.Agent().ServiceDeregister(serviceUUID); err != nil {
			fmt.Printf("[Gorra] Deregister Service %s Failed: %v\n", serviceUUID, err)
		}

		fmt.Printf("[Gorra] Deregister Service %s Success\n", serviceUUID)

		os.Exit(200)
	}()

	fmt.Printf("[Gorra] Service %s is up and running\n", serviceUUID)

	return nil
}
