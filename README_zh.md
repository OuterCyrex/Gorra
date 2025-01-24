# Gorra - Golang 微服务脚手架

[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Stars](https://img.shields.io/badge/github-stars-blue)](https://github.com/OuterCyrex/Gorra/stargazers)

**Gorra** 是一个基于 Golang 的微服务脚手架，旨在简化微服务的开发流程。通过集成 Consul 和 Nacos，它提供了便捷的服务注册、配置管理以及健康检查功能。开发者可以通过简单的配置和代码，快速搭建高性能、可扩展的微服务架构。

## 特性

- **服务注册与发现**：支持 Consul 和 Nacos，无缝集成服务注册与发现。
- **配置管理**：动态加载配置文件，支持热更新。
- **健康检查**：内置健康检查机制，确保服务的高可用性。
- **优雅启动与关闭**：支持服务的优雅启动和关闭。
- **简单易用**：简洁的 API 设计，快速上手。

## 安装

```bash
go get github.com/OuterCyrex/Gorra
```

## 快速开始

### 1. 配置文件

创建一个 `config.yaml` 文件，配置服务的基本信息和注册中心的连接信息。例如：

```yaml
# config.yaml
host: '127.0.0.1'
port: 8848
namespace: 'example_namespace'
dataId: 'user_srv'
group: 'Debug'
```

Nacos 中的配置文件应当按如下字段进行配置：

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

### 2. 初始化服务

在你的项目中，使用 `Gorra` 初始化服务并启动：

```go
package main

import (
	"github.com/OuterCyrex/Gorra"
	proto "user_srv/proto/.UserProto"
	"google.golang.org/grpc"
)

func main() {
    // 从Nacos中初始化您的服务配置
	serverConfig, err := Gorra.InitConfig("example_config.yaml")
	if err != nil {
		panic(err)
	}

    // 注册grpc服务
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &UserServer{})

    // 运行服务并将服务注册至Consul
	err = Gorra.ServerRun(server, serverConfig)
	if err != nil {
		panic(err)
	}
}
```

## 文档

- **配置文件说明**：`config.yaml` 文件用于配置服务的基本信息。
- **服务注册**：目前仅支持 Consul 服务注册。
- **健康检查**：内置健康检查机制，确保服务的高可用性。
- **优雅关闭**：支持服务的优雅关闭，确保资源正确释放。

## 示例项目

- 代办

## 许可

Gorra 采用 [MIT 许可证](https://opensource.org/licenses/MIT)。