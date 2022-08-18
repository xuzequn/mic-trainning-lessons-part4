package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"math/rand"
	"time"
)

const resName = "cart-order"

func main() {
	err := sentinel.InitDefault()
	if err != nil {
		panic(err)
	}
	var all int
	var through int
	var block int
	ch := make(chan struct{})
	var rules []*flow.Rule
	rule := &flow.Rule{
		Resource:               resName,
		TokenCalculateStrategy: flow.WarmUp,
		ControlBehavior:        flow.Reject,
		Threshold:              1000, // 与下面的参数一起决定流量控制的灵敏度，
		WarmUpPeriodSec:        30,   // 在这个设定的周期内允许的最大请求数量为上面参数的值。以秒为单位
	}
	rules = append(rules, rule)
	_, err = flow.LoadRules(rules)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		go func() {
			for {
				all++
				entry, blockError := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
				if blockError != nil {
					block++
					//fmt.Println("流量太大，开启限流")
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					through++
					//fmt.Println("限流通过")
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					entry.Exit()
				}
			}
		}()
	}
	go func() {
		var oldAll int
		var oldThough int
		var oldBlock int
		for {
			a := all - oldAll
			oldAll = all

			t := through - oldThough
			oldThough = through

			b := block - oldBlock
			oldBlock = block
			time.Sleep(time.Second * 1)
			fmt.Printf("all:%d,though:%d,block:%d \n", a, t, b)
		}
	}()
	<-ch
}
