package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver" // 让 grpc 可以解析consul协议
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"mic-trainning-lesson-part4/cartorder_web/handler"
	"mic-trainning-lesson-part4/cartorder_web/middleware"
	"mic-trainning-lesson-part4/cartorder_web/order"
	"mic-trainning-lesson-part4/internal"
	"mic-trainning-lesson-part4/internal/register"
	"mic-trainning-lesson-part4/util"
	"os"
	"os/signal"
	"syscall"
)

var (
	addr           string
	consulRegistry register.ConsulRegistry
	randomId       string
)

func init() {
	randomPort := util.GenRandomPort()

	if !internal.AppConf.Debug {
		internal.AppConf.CartOrderWebConfig.Port = randomPort
	}
	randomId = uuid.NewV4().String()
	consulRegistry = register.NewConsulRegistry(internal.AppConf.ConsulConfig.Host, int(internal.AppConf.ConsulConfig.Port))
	consulRegistry.Register(internal.AppConf.CartOrderWebConfig.SrvName, randomId,
		internal.AppConf.CartOrderWebConfig.Port, internal.AppConf.CartOrderWebConfig.SrvType,
		internal.AppConf.CartOrderWebConfig.Tags)
	addr = fmt.Sprintf("%s:%d", internal.AppConf.CartOrderWebConfig.Host, internal.AppConf.CartOrderWebConfig.Port)

}

func main() {

	r := gin.Default()
	CartOrderGroup := r.Group("/v1/cart").Use(middleware.Tracing())
	{
		CartOrderGroup.GET("/list/:accountId", handler.ShopCartListHandler)
		CartOrderGroup.POST("/add", handler.AddHandler)
		CartOrderGroup.POST("/update", handler.UpdateHandler)
		CartOrderGroup.POST("/delete", handler.DelHandler)
	}
	orderGroup := r.Group("/v1/order").Use(middleware.Tracing())
	{
		orderGroup.GET("", order.ListHandler)
		orderGroup.GET("/:id", order.Detail)
		orderGroup.GET("/add", order.CreateOrder)

	}
	r.GET("/health", handler.HealthHandler)

	go func() {
		err := r.Run(addr)
		if err != nil {
			zap.S().Panic(addr + "启动失败" + err.Error())
		} else {
			zap.S().Info(addr + "启动成功")
		}
	}()
	q := make(chan os.Signal)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-q
	err := consulRegistry.DeRegister(randomId)
	if err != nil {
		zap.S().Panic("注销失败" + randomId + "：" + err.Error())
	} else {
		zap.S().Info("注销成功" + randomId)
	}
}
