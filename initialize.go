package Gorra

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

func InitConfig(path string) (ServerConfig, error) {
	nacosConfig := getNacosConfig(path)
	sc := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.Host,
			Port:   nacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		return ServerConfig{}, err
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
	})
	if err != nil {
		return ServerConfig{}, err
	}

	serverConfig := ServerConfig{}

	err = json.Unmarshal([]byte(content), &serverConfig)
	if err != nil {
		return ServerConfig{}, err
	}

	fmt.Println("[Gorra] Initialize Server Config Success")

	return serverConfig, nil
}

func getNacosConfig(path string) NacosConfig {
	YAMLFile := path
	v := viper.New()
	v.SetConfigFile(YAMLFile)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	var nacos NacosConfig

	if err := v.Unmarshal(&nacos); err != nil {
		panic(err)
	}

	fmt.Printf("[Gorra] Get Nacos YMAL Config Success\n")

	return nacos
}
