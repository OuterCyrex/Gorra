package GorraAPI

type RegistryInfo struct {
	Name    string       `mapstructure:"name" json:"name"`
	Address string       `mapstructure:"address" json:"address"`
	Port    int          `mapstructure:"port" json:"port"`
	Tags    []string     `mapstructure:"tags" json:"tags"`
	Consul  ConsulConfig `mapstructure:"consul" json:"consul"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstruct:"host"`
	Port      uint64 `mapstruct:"port"`
	Namespace string `mapstruct:"namespace"`
	DataId    string `mapstruct:"dataId"`
	Group     string `mapstruct:"group"`
}

type BaseConfig interface {
	GetRegistryInfo() RegistryInfo
}
