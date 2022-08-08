package order

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mic-trainning-lesson-part4/internal"
	"mic-trainning-lesson-part4/proto/pb"
	"net/http"
)

func ListHandler(c *gin.Context) {
	//accountId, _ := c.Get("accountId")
	//claims, _ := c.Get("claims")
	//pageNoStr := c.DefaultQuery("pageNo", "0")
	//pageNo, _ := strconv.Atoi(pageNoStr)
	//
	//pageSizeStr := c.DefaultQuery("pageSize", "0")
	//pageSize, _ := strconv.Atoi(pageSizeStr)
	//
	//reqPb := &pb.OrderPagingReq{
	//	PageNo:   int32(pageNo),
	//	PageSize: int32(pageSize),
	//}
	//
	//// 如果是管理员用户则返回所有订单
	//customClaims := claims.(*jwt_op.CustromClaims)
	//if customClaims.AuthorityId == 1 {
	//	reqPb.AccountId = int32(accountId.(uint))
	//}
	ctx := context.WithValue(context.Background(), "webContext", c)
	res, err := internal.OrderClient.OrderList(ctx, &pb.OrderPagingReq{
		PageNo:   2,
		PageSize: 1,
	})
	if err != nil {
		zap.S().Error(err)
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      "",
		"total":    res.Total,
		"itemList": res.ItemList,
	})
}

func Detail(c *gin.Context) {
	ctx := context.WithValue(context.Background(), "webContext", c)
	res, err := internal.OrderClient.OrderDetail(ctx, &pb.OrderItemReq{
		Id: 7,
	})
	if err != nil {
		zap.S().Error(err)
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      "",
		"itemList": res.ProductList,
	})
}

func CreateOrder(c *gin.Context) {
	ctx := context.WithValue(context.Background(), "webContext", c)
	res, err := internal.OrderClient.CreateOrder(ctx, &pb.OrderItemReq{
		AccountId: 1,
		Addr:      "北京",
		PostCode:  "10010",
		Receiver:  "xuzequn",
		Mobile:    "13500000000",
	})
	if err != nil {
		zap.S().Error(err)
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":     "",
		"orderNo": res.OrderNum,
	})
}
