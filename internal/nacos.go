package internal

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	NameSpace string `mapstructure:"namespace"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
