package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"math/rand"
	"time"
)

const resName = "cart-order"

func main() {
	conf := config.NewDefaultConfig()
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		panic(err)
	}
	var rules []*flow.Rule
	rule := &flow.Rule{
		Resource:               resName,
		TokenCalculateStrategy: flow.Direct,
		ControlBehavior:        flow.Reject,
		Threshold:              10,   // 与下面的参数一起决定流量控制的灵敏度，
		StatIntervalInMs:       1000, // 在这个设定的周期内允许的最大请求数量为上面参数的值。如果是1000，就是等于QPS数。
	}
	rules = append(rules, rule)
	_, err = flow.LoadRules(rules)
	if err != nil {
		panic(err)
	}

	ch := make(chan struct{})
	for i := 0; i < 2; i++ {
		go func() {
			for {
				entry, blockError := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
				if blockError != nil {
					fmt.Println("流量太大，开启限流")
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					fmt.Println("限流通过")
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					entry.Exit()
				}
			}
		}()
	}
	go func() {
		time.Sleep(3 * time.Second)
		_, err = flow.LoadRules([]*flow.Rule{
			{
				Resource:               resName,
				TokenCalculateStrategy: flow.Direct,
				ControlBehavior:        flow.Reject,
				Threshold:              180,
				StatIntervalInMs:       1000,
			},
		})
		if err != nil {
			panic(err)
		}
	}()
	<-ch
}
