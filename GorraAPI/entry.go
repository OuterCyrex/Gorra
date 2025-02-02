package GorraAPI

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"os"
	"os/signal"
	"syscall"
)

func RunRouter(router *gin.Engine, config BaseConfig) error {
	registryId := fmt.Sprintf("%s", uuid.NewV4())

	registryClient, err := newRegistryClient(config.GetRegistryInfo().Consul.Host, config.GetRegistryInfo().Consul.Port)
	if err != nil {
		return err
	}
	err = registryClient.register(
		config.GetRegistryInfo().Address,
		config.GetRegistryInfo().Port,
		config.GetRegistryInfo().Name,
		config.GetRegistryInfo().Tags,
		fmt.Sprintf("%s", registryId),
	)

	if err != nil {
		return err
	}

	fmt.Printf("\u001B[31m[GorraAPI] API Server Runs On Port %d......\u001B[0m\n", config.GetRegistryInfo().Port)

	go func() {
		_ = router.Run(fmt.Sprintf(":%d", config.GetRegistryInfo().Port))
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = registryClient.deRegister(registryId)
	if err == nil {
		fmt.Printf("[Gorra]: API Gateway Deregistry %v Success\n", registryId)
	}

	return nil
}

type registry struct {
	Host   string
	Port   int
	Client *api.Client
}

type registryClient interface {
	register(address string, port int, name string, tags []string, id string) error
	deRegister(serviceId string) error
}

func newRegistryClient(host string, port int) (registryClient, error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", host, port)
	client, err := api.NewClient(cfg)
	if err != nil {
		msg := fmt.Sprintf("[Gorra]: API Register Failed: %v", err)
		return nil, errors.New(msg)
	}
	return &registry{
		Host:   host,
		Port:   port,
		Client: client,
	}, nil
}

func (r *registry) register(address string, port int, name string, tags []string, id string) error {

	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "15s",
	}

	registration := &api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Port:    port,
		Tags:    tags,
		Address: address,
		Check:   check,
	}

	err := r.Client.Agent().ServiceRegister(registration)
	return err
}

func (r *registry) deRegister(serviceId string) error {
	err := r.Client.Agent().ServiceDeregister(serviceId)
	return err
}
