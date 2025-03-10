# Gorra - Golang Microservice Scaffold

- [中文文档](README_zh.md)

[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Stars](https://img.shields.io/badge/github-stars-blue)](https://github.com/OuterCyrex/Gorra/stargazers)

**Gorra** is a Golang microservice scaffold designed to simplify the development process of microservices. By integrating Consul and Nacos, it provides convenient service registration, configuration management, and health check capabilities. Developers can quickly build high-performance, scalable microservice architectures with simple configurations and code.

## Features

- **Service Registration and Discovery**: Supports Consul and Nacos for seamless service registration and discovery.
- **Configuration Management**: Dynamically loads configuration files with hot-reloading support.
- **Health Check**: Built-in health check mechanism to ensure high availability of services.
- **Graceful Start and Shutdown**: Supports graceful start and shutdown of services.
- **Easy to Use**: Simple API design for quick adoption.

## Installation

```bash
go get github.com/OuterCyrex/Gorra
```

## Quick Start

### 1. Configuration File

Create a `config.yaml` file to configure the service details and registry connection. For example:

```yaml
# config.yaml
host: '127.0.0.1'
port: 8848
namespace: 'example_namespace'
dataId: 'user_srv'
group: 'Debug'
```

In Server Part, namely `GorraSrv`, The configuration file in Nacos should be configured according to the following fields:

```json
{
  "name": "user_srv",
  "addr":"127.0.0.1",
  "tags": ["user", "grpc", "service"],
  "mysql": {
    "host": "127.0.0.1",
    "port": 3306,
    "username": "your_username",
    "password": "your_password",
    "db": "your_db"
  },
  "consul": {
    "host": "127.0.0.1",
    "port": 8500
  }
}
```

But In the API Gateway part, the configuration file in Nacos can be set by developers according to their own needs, as long as it implements the `GorraAPI.BaseConfig` interface.

```go
// Implement the GetRegistryInfo method of BaseConfig
func (m APIConfig) GetRegistryInfo() GorraAPI.RegistryInfo {
    return GorraAPI.RegistryInfo{
        Name:    m.Name,
        Address: m.Address,
        Port:    m.Port,
        Tags:    m.Tags,
        Consul: GorraAPI.ConsulConfig{
            Host: m.Consul.Host,
            Port: m.Consul.Port,
        },
    }
}
```

### 2. Initialize the Service

In your project, initialize and start the service using `GorraSrv`:

```go
package main

import (
	"github.com/OuterCyrex/Gorra/GorraSrv"
	proto "user_srv/proto/.UserProto"
	"google.golang.org/grpc"
)

func main() {
	// pull Config from Nacos and initialize it
	serverConfig, err := GorraSrv.InitConfig("example_config.yaml")
	if err != nil {
		panic(err)
	}

	// register grpc server
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &UserServer{})

	// run the server and register it in Consul
	err = GorraSrv.ServerRun(server, serverConfig)
	if err != nil {
		panic(err)
	}
}
```

### 3.Initialize the API Gateway

After initialize the, initialize and start the service using `GorraAPI`:

```go
package main

import (
    "github.com/OuterCyrex/Gorra/GorraAPI"
)

type APIConfig struct {
    Name     string         `mapstructure:"name" json:"name"`
    Address  string         `mapstructure:"address" json:"address"`
    Port     int            `mapstructure:"port" json:"port"`
    Tags     []string       `mapstructure:"tags" json:"tags"`
    JwtKey   string         `mapstructure:"jwtKey" json:"jwtKey"`
    Consul   ConsulConfig   `mapstructure:"consul" json:"consul"`
    GoodsSrv GoodsSrvConfig `mapstructure:"goodsSrv" json:"goodsSrv"`
}

type ConsulConfig struct {
    Host string `mapstructure:"host" json:"host"`
    Port int    `mapstructure:"port" json:"port"`
}

type GoodsSrvConfig struct {
    Name string `mapstructure:"name" json:"name"`
}

// Implement the GetRegistryInfo method of BaseConfig
func (m APIConfig) GetRegistryInfo() GorraAPI.RegistryInfo {
    return GorraAPI.RegistryInfo{
       Name:    m.Name,
       Address: m.Address,
       Port:    m.Port,
       Tags:    m.Tags,
       Consul: GorraAPI.ConsulConfig{
          Host: m.Consul.Host,
          Port: m.Consul.Port,
       },
    }
}

func main() {
    c := &APIConfig{}

    // pull Config from Nacos
    cf, err := GorraAPI.InitConfig("test/config.yaml", c)

    if err != nil {
       panic(err)
    }

    // initialize the Router
    r := GorraAPI.KeepAliveRouters("v1")

    // run the router serve
    err = GorraAPI.RunRouter(r, cf)

    if err != nil {
       panic(err)
    }
}
```

## Documentation

- **Configuration File**: The `config.yaml` file is used to configure service details.
- **Service Registration**: Now supports Consul as the services registery center.
- **Health Check**: Built-in health check mechanism to ensure high availability.
- **Graceful Shutdown**: Supports graceful shutdown to ensure proper resource release.

## Example Projects

- TODO

## License

Gorra is licensed under the [MIT License](https://opensource.org/licenses/MIT).