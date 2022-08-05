package biz

import (
	"context"
	"fmt"
	"mic-trainning-lesson-part4/internal"
	"mic-trainning-lesson-part4/proto/pb"
	"testing"
)

func TestCartOrderServer_CreateOrder(t *testing.T) {
	res, err := internal.OrderClient.CreateOrder(context.Background(), &pb.OrderItemReq{
		AccountId: 1,
		Addr:      "北京",
		PostCode:  "10010",
		Receiver:  "xuzequn",
		Mobile:    "13500000000",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
