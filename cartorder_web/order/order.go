package order

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mic-trainning-lesson-part4/cartorder_web/req"
	"mic-trainning-lesson-part4/internal"
	"mic-trainning-lesson-part4/proto/pb"
	"net/http"
	"strconv"
)

func ListHandler(c *gin.Context) {

	accountIdStr := c.DefaultQuery("accountId", "0")
	accountId, _ := strconv.Atoi(accountIdStr)

	if accountId < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
	}

	pageNoStr := c.DefaultQuery("pageNo", "0")
	pageNo, _ := strconv.Atoi(pageNoStr)

	pageSizeStr := c.DefaultQuery("pageSize", "0")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	reqPb := &pb.OrderPagingReq{
		PageNo:    int32(pageNo),
		PageSize:  int32(pageSize),
		AccountId: int32(accountId),
	}

	// 如果是管理员用户则返回所有订单

	ctx := context.WithValue(context.Background(), "webContext", c)
	res, err := internal.OrderClient.OrderList(ctx, reqPb)
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
	orderReq := req.OrderReq{}
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部错误" + err.Error(),
		})
	}
	ctx := context.WithValue(context.Background(), "webContext", c)
	orderItemPb := ConventOrderModel2Pb(orderReq)
	res, err := internal.OrderClient.CreateOrder(ctx, orderItemPb)
	if err != nil {
		zap.S().Error(err)
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":     "",
		"orderNo": res.OrderNum,
	})
}

func ConventOrderModel2Pb(req req.OrderReq) *pb.OrderItemReq {
	return &pb.OrderItemReq{
		Id:        req.Id,
		AccountId: req.AccountId,
		Addr:      req.Addr,
		PostCode:  req.PostCode,
		Receiver:  req.Receiver,
		Mobile:    req.ReceiverMobile,
		PayType:   req.PayType,
	}
}
