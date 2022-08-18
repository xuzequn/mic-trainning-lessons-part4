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
	var rules []*flow.Rule
	rule := &flow.Rule{
		Resource:               resName,
		TokenCalculateStrategy: flow.Direct,
		ControlBehavior:        flow.Throttling,
		Threshold:              10,
		StatIntervalInMs:       1000,
	}
	rules = append(rules, rule)
	_, err = flow.LoadRules(rules)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 20; i++ {
		entry, blockError := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
		if blockError != nil {
			fmt.Println("流量太大，开启限流")
			time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
		} else {
			fmt.Println("限流通过")
			time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
			entry.Exit()
		}
		time.Sleep(150 * time.Millisecond)
	}
}
