package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

const groupName = "testGroup"

func GetMqAddr() string {
	//mqAddr := fmt.Sprintf("%s:%d", internal.AppConf.RocketMQConfig.Host,
	//	internal.AppConf.RocketMQConfig.Port)
	mqAddr := "127.0.0.1:9876"
	return mqAddr

}

func ProduceMsg(mqAddr string, topic string) {
	p, err := rocketmq.NewProducer( // 普通消息生产者
		producer.WithGroupName(groupName),
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{mqAddr})),
		producer.WithRetry(2),
	)
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		zap.S().Error("生产者错误" + err.Error())
		os.Exit(1)
	}
	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("Hello Happy Mall" + strconv.Itoa(i)),
		}
		msg.WithDelayTimeLevel(3)
		r, err := p.SendSync(context.Background(), msg)
		if err != nil {
			zap.S().Error("发送消息错误" + err.Error())
		} else {
			zap.S().Info("发送消息成功" + r.String() + "-" + r.MsgID)
		}
	}
	err = p.Shutdown()
	if err != nil {
		zap.S().Error("生产者shutdown" + err.Error())
		os.Exit(1)
	}
}

func ComsumeMsg(mqAddr string, topic string) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(groupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{mqAddr})),
	)
	if err != nil {
		panic(err)
	}
	err = c.Subscribe(topic, consumer.MessageSelector{},
		func(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgList {
				fmt.Printf("订阅消息，消费%v \n", msgList[i])
			}
			return consumer.ConsumeSuccess, nil
		})
	if err != nil {
		zap.S().Error("消费消息错误" + err.Error())
	}
	err = c.Start()
	if err != nil {
		zap.S().Error("开启消费这错误" + err.Error())
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		zap.S().Error("shutdown消费者错误" + err.Error())
	}

}

func main() {
	topic := "HappyMall"
	mqAddr := GetMqAddr()
	ProduceMsg(mqAddr, topic)
	ComsumeMsg(mqAddr, topic)
}
