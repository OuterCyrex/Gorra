package GorraAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"reflect"
)

func InitConfig(YamlPath string, configStruct BaseConfig) (BaseConfig, error) {
	nacosConfig := getNacosConfig(YamlPath)
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
		msg := fmt.Sprintf("[Gorra]: Create Nacos Client Failed: %v", err)
		return nil, errors.New(msg)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
	})
	if err != nil {
		msg := fmt.Sprintf("[Gorra]: Get Nacos JSON Failed: %v", err)
		return nil, errors.New(msg)
	}

	configPtr := reflect.ValueOf(configStruct)
	if configPtr.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("configStruct must be a pointer")
	}
	configElem := configPtr.Elem()

	err = json.Unmarshal([]byte(content), configElem.Addr().Interface())
	if err != nil {
		msg := fmt.Sprintf("[Gorra]: Unmarshal Nacos JSON Failed: %v", err)
		return nil, errors.New(msg)
	}

	return configElem.Interface().(BaseConfig), nil
}

func getNacosConfig(YamlPath string) NacosConfig {
	YAMLFile := YamlPath
	v := viper.New()
	v.SetConfigFile(YAMLFile)
	if err := v.ReadInConfig(); err != nil {
		panic("[Gorra]: Viper Read YAMLFile failed")
	}

	var nacos NacosConfig

	if err := v.Unmarshal(&nacos); err != nil {
		panic("[Gorra]: Viper UnMarshal YAMLFile failed")
	}

	return nacos
}
