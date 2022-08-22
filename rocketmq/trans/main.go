package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"time"
)

type HappyListener struct {
}

func (hl HappyListener) ExecuteLocalTransaction(*primitive.Message) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

func (hl HappyListener) CheckLocalTransaction(*primitive.MessageExt) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

func main() {
	mqAddr := "127.0.0.1:9876"
	p, err := rocketmq.NewTransactionProducer( // 开启事物消息生产者
		HappyListener{},
		producer.WithNameServer([]string{mqAddr}),
	)
	if err != nil {
		panic(err) // 生产环境禁用panic
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}
	res, err := p.SendMessageInTransaction(context.Background(),
		primitive.NewMessage("HapplyTransactionTopic", []byte("从0到GO语言为服务架构师")))
	fmt.Println(res.Status)
	if err != nil {
		panic(err)
	}
	fmt.Printf("发送成功")
	time.Sleep(time.Second * 3600)
	err = p.Shutdown()
	if err != nil {
		panic(err)
	}
}
