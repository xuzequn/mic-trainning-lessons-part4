package internal

import (
	"encoding/json"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mic-trainning-lesson-part4/proto/pb"
	//"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"mic-trainning-lesson-part4/util/otgrpc"
)

var AppConf AppConfig
var NacosConf NacosConfig

var ShopCartClient pb.ShopCartServiceClient

var OrderClient pb.OrderServiceClient
var ProductClient pb.ProductServiceClient
var StockClient pb.StockServiceClient

//var ViperConf ViperConfig
var fileName = "dev-config.yaml"

//var fileName = "../../dev-config.yaml"

func initNacos() {
	v := viper.New()
	v.SetConfigFile(fileName)
	v.ReadInConfig()
	err := v.Unmarshal(&NacosConf)
	fmt.Println(NacosConf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func initFromNacos() {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: NacosConf.Host,
			Port:   NacosConf.Port,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         NacosConf.NameSpace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "nacos/log",
		CacheDir:            "nacos/cache",
		LogLevel:            "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: NacosConf.DataId,
		Group:  NacosConf.Group,
	})
	if err != nil {
		panic(err)
	}
	if content == "" {
		panic("配置文件为空")
	}
	fmt.Println(content)
	json.Unmarshal([]byte(content), &AppConf)
	fmt.Println(AppConf)
}

func init() {
	initNacos()
	initFromNacos()
	fmt.Println("config初始化完成。。。")
	InitRedis()
	InitDB()

	addr := fmt.Sprintf("%s:%d", AppConf.ConsulConfig.Host, AppConf.ConsulConfig.Port)
	dialAddr := fmt.Sprintf("consul://%s/%s?wait=14s", addr, AppConf.CartOrderSrvConfig.SrvName)
	conn, err := grpc.Dial(dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robbin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal(err)
		panic(err)
	}
	ShopCartClient = pb.NewShopCartServiceClient(conn)
	OrderClient = pb.NewOrderServiceClient(conn)

	productSrvAddr := fmt.Sprintf("consul://%s/%s?wait=14s", addr, AppConf.ProductSrvConfig.SrvName)
	productConn, err := grpc.Dial(productSrvAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robbin"}`),
		// grpc 一元连接器           otgrpc.opentracing 客户端拦截器      opentracing 的全局追踪器
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal(err)
		panic(err)
	}
	ProductClient = pb.NewProductServiceClient(productConn)

	stockSrvAddr := fmt.Sprintf("consul://%s/%s?wait=14s", addr, AppConf.StockSrvConfig.SrvName)
	stockConn, err := grpc.Dial(stockSrvAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robbin"}`),
	)
	if err != nil {
		zap.S().Fatal(err)
		panic(err)
	}

	StockClient = pb.NewStockServiceClient(stockConn)
}
