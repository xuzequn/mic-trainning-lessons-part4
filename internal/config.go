package internal

type CartOrderSrvConfig struct {
	SrvName string   `mapstructure:"srvName" json:"srvName"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
	SrvType string   `mapstructure:"srvType" json:"srvType"`
}

type CartOrderWebConfig struct {
	SrvName string   `mapstructure:"srvName" json:"srvName"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
	SrvType string   `mapstructure:"srvType" json:"srvType"`
}

type ProductSrvConfig struct {
	SrvName string   `mapstructure:"srvName" json:"srvName"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
	SrvType string   `mapstructure:"srvType" json:"srvType"`
}

type StockSrvConfig struct {
	SrvName string   `mapstructure:"srvName" json:"srvName"`
	Host    string   `mapstructure:"host" json:"host"`
	Port    int      `mapstructure:"port" json:"port"`
	Tags    []string `mapstructure:"tags" json:"tags"`
	SrvType string   `mapstructure:"srvType" json:"srvType"`
}

type AppConfig struct {
	DBConfig           DBConfig           `mapstructure:"db" json:"db"`
	RedisConfig        RedisConfig        `mapstructure:"redis" json:"redis"`
	ConsulConfig       ConsulConfig       `mapstructure:"consul" json:"consul"`
	CartOrderSrvConfig CartOrderSrvConfig `mapstructure:"cart_order_srv" json:"cart_order_srv"`
	CartOrderWebConfig CartOrderWebConfig `mapstructure:"cart_order_web" json:"cart_order_web"`
	ProductSrvConfig   ProductSrvConfig   `mapstructure:"product_srv" json:"product_srv"`
	StockSrvConfig     StockSrvConfig     `mapstructure:"stock_srv" json:"stock_srv"`
	JWTConfig          JWTConfig          `mapstructure:"jwt" json:"jwt"`
	Debug              bool               `mapstructure:"debug" json:"debug"`
}
