package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mic-trainning-lesson-part4/cartorder_srv/biz"
	"mic-trainning-lesson-part4/internal"
	"mic-trainning-lesson-part4/internal/register"
	"mic-trainning-lesson-part4/proto/pb"
	"mic-trainning-lesson-part4/util"
	"mic-trainning-lesson-part4/util/otgrpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	consulRegistry register.ConsulRegistry
	randomId       string
)

func init() {
	randomPort := util.GenRandomPort()
	if !internal.AppConf.Debug {
		internal.AppConf.CartOrderSrvConfig.Port = randomPort
	}
	randomId = uuid.NewV4().String()
	consulRegistry = register.NewConsulRegistry(internal.AppConf.ConsulConfig.Host,
		int(internal.AppConf.ConsulConfig.Port))
}

func main() {
	/*
			1、生成proto对应的文件
			2、简历biz目录，生成对应接口。
		    3、拷贝之前main文件的函数、查缺补漏
	*/

	//port := util.GenRandomPort()
	port := internal.AppConf.CartOrderSrvConfig.Port
	addr := fmt.Sprintf("%s:%d", internal.AppConf.CartOrderSrvConfig.Host, port)

	// 链路追踪
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", internal.AppConf.JaegerConfig.AgentHost, internal.AppConf.JaegerConfig.AgentPort),
		},
		ServiceName: "xzqMall",
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	defer closer.Close()
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)

	// 将定义的对象注册grpc服务
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	pb.RegisterShopCartServiceServer(server, &biz.CartOrderServer{})
	pb.RegisterOrderServiceServer(server, &biz.CartOrderServer{})
	// 启动服务监听
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Error("cartorder_srv 启动异常" + err.Error())
		panic(err)
	}
	// grpc 服务的健康检查  注册服务健康检查  启动的grpc  健康检查方法
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	// 注册服务
	err = consulRegistry.Register(internal.AppConf.CartOrderSrvConfig.SrvName, randomId,
		internal.AppConf.CartOrderSrvConfig.Port, internal.AppConf.CartOrderSrvConfig.SrvType, internal.AppConf.CartOrderSrvConfig.Tags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("%s启动在%d", randomId, port))

	mqAddr := "127.0.0.1:9876"
	pushConsumer, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{mqAddr}),
		consumer.WithGroupName("HappyOrderTimeOut"),
	)
	pushConsumer.Subscribe("timeout_order_info", consumer.MessageSelector{},
		func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		})
	go func() {
		err = server.Serve(listen)
		if err != nil {
			zap.S().Panic(addr + "启动失败" + err.Error())
			fmt.Println(addr + "启动失败" + err.Error())
		} else {
			zap.S().Info(addr + "启动成功")
		}
	}()
	q := make(chan os.Signal)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q
	err = consulRegistry.DeRegister(randomId)
	if err != nil {
		zap.S().Panic("注销失败" + randomId + ":" + err.Error())
	} else {
		zap.S().Info("注销成功" + randomId)
	}
}
