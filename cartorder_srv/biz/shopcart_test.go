package biz

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // 让 grpc 可以解析consul协议
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mic-trainning-lesson-part4/internal"
	"mic-trainning-lesson-part4/proto/pb"
	"testing"
)

var (
	client pb.ShopCartServiceClient
)

func init() {
	addr := fmt.Sprintf("%s:%d", internal.AppConf.ConsulConfig.Host, internal.AppConf.ConsulConfig.Port)
	dialAddr := fmt.Sprintf("consul://%s/%s?wait=14s", addr, internal.AppConf.CartOrderSrvConfig.SrvName)
	conn, err := grpc.Dial(dialAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robbin"}`),
	)
	if err != nil {
		zap.S().Fatal(err)
		panic(err)
	}
	client = pb.NewShopCartServiceClient(conn)
}

func TestShopCartServer_ShopCartItemList(t *testing.T) {
	res, err := client.ShopCartItemList(context.Background(), &pb.AccountReq{AccountId: 1})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestShopCartServer_AddShopCartItem(t *testing.T) {
	//res, err := client.AddShopCartItem(context.Background(), &pb.ShopCartReq{
	//	AccountId: 1,
	//	ProductId: 6,
	//	Num:       1,
	//	Checked:   true,
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(res)

	for i := 0; i < 5; i++ {
		res, err := client.AddShopCartItem(context.Background(), &pb.ShopCartReq{
			AccountId: 1,
			ProductId: 6 + int32(i),
			Num:       1,
			Checked:   true,
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(res)
	}
}

func TestShopCartServer_UpdateShopCartItem(t *testing.T) {
	res, err := client.UpdateShopCartItem(context.Background(), &pb.ShopCartReq{
		AccountId: 1,
		ProductId: 7,
		Num:       2,
		Checked:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)

}

func TestShopCartServer_DeleteShopCartItem(t *testing.T) {
	r, err := client.DeleteShopCartItem(context.Background(), &pb.DelShopCartItem{
		AccountId: 1,
		ProductId: 8,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r)
}
