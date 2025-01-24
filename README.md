# Gorra - Golang Microservice Scaffold

- [中文](README_zh.md)

[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
![Stargazers](https://github.com/OuterCyrex/Gorra/stargazers)

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

yaml复制

```yaml
# config.yaml
host: '127.0.0.1'
port: 8848
namespace: 'example_namespace'
dataId: 'user_srv'
group: 'Debug'
```

### 2. Initialize the Service

In your project, initialize and start the service using `Gorra`:

go复制

```go
package main

import (
	"github.com/OuterCyrex/Gorra"
	proto "user_srv/proto/.UserProto"
	"google.golang.org/grpc"
)

func main() {
    // pull Config from Nacos and initialize it
	serverConfig, err := Gorra.InitConfig("example_config.yaml")
	if err != nil {
		panic(err)
	}

    // register grpc server
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &UserServer{})

    // run the server and register it in Consul
	err = Gorra.ServerRun(server, serverConfig)
	if err != nil {
		panic(err)
	}
}
```

## Documentation

- **Configuration File**: The `config.yaml` file is used to configure service details and registry connections. For more details, see the [Configuration Documentation](https://kimi.moonshot.cn/chat/docs/config.md).
- **Service Registration**: Supports Consul and Nacos, specified via the configuration file.
- **Health Check**: Built-in health check mechanism to ensure high availability.
- **Graceful Shutdown**: Supports graceful shutdown to ensure proper resource release.

## Example Projects

- [User Service Example](https://kimi.moonshot.cn/chat/examples/user_srv): A complete example of a user service demonstrating how to build a microservice with Gorra.

## Contributing

PRs and issue reports are welcome! For more details, see the [Contributing Guide](https://kimi.moonshot.cn/chat/CONTRIBUTING.md).

## License

Gorra is licensed under the [MIT License](https://kimi.moonshot.cn/chat/LICENSE).