package handler

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

func ShopCartListHandler(c *gin.Context) {
	accountStr := c.Param("accountId")
	accountId, err := strconv.Atoi(accountStr)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数错误",
		})
		return
	}
	res, err := internal.ShopCartClient.ShopCartItemList(context.Background(), &pb.AccountReq{AccountId: int32(accountId)})
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

func AddHandler(c *gin.Context) {
	var shopCartReq req.ShopCartReq
	err := c.ShouldBind(&shopCartReq)
	if err != nil {
		zap.S().Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数解析错误",
		})
		return
	}
	r := ConvertShopCartReq2Pb(shopCartReq)
	res, err := internal.ShopCartClient.AddShopCartItem(context.Background(), r)
	if err != nil {
		zap.S().Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "添加购物车失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": res,
	})
}

func DelHandler(c *gin.Context) {
	var cartReq pb.ShopCartReq
	err := c.ShouldBindJSON(&cartReq)
	if err != nil {
		zap.S().Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数解析错误",
		})
		return
	}
	_, err = internal.ShopCartClient.DeleteShopCartItem(context.Background(),
		&pb.DelShopCartItem{AccountId: int32(cartReq.AccountId), ProductId: int32(cartReq.ProductId)})
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "删除购物车失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "",
	})
}

func UpdateHandler(c *gin.Context) {
	var cartReq req.ShopCartReq
	err := c.ShouldBind(&cartReq)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "参数解析错误",
		})
		return
	}
	r := ConvertShopCartReq2Pb(cartReq)
	_, err = internal.ShopCartClient.UpdateShopCartItem(context.Background(), r)
	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": "更新产品失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "",
	})

}

func ConvertShopCartReq2Pb(shopcartreq req.ShopCartReq) *pb.ShopCartReq {
	item := pb.ShopCartReq{
		AccountId: shopcartreq.AccountId,
		ProductId: shopcartreq.ProductId,
		Num:       shopcartreq.Num,
		Checked:   *shopcartreq.Checked,
	}
	if shopcartreq.Id != 0 {
		item.Id = shopcartreq.Id
	}
	return &item
}
