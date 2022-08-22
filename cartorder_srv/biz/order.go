package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"mic-trainning-lesson-part4/cartorder_srv/model"
	"mic-trainning-lesson-part4/custom_error"
	"mic-trainning-lesson-part4/internal"
	"mic-trainning-lesson-part4/proto/pb"
)

type OrderListener struct {
	Id          int32
	Detail      string
	OrderNo     string
	OrderAmount float32
	AccountId   int32
	Status      codes.Code
	Addr        string
	Receiver    string
	Mobile      string
	PostCode    string
}

func (ol *OrderListener) ExecuteLocalTransaction(*primitive.Message) primitive.LocalTransactionState {
	// 1 半消息
	// 2、 执行库存扣减
	// 3、 返回成功
	// 4、 执行本地消息

	/*
		1、拿到购物车内的选中商品
		2、订单总金额， 不要相信前端数据
		3、扣减库存 stock_srv
		4、把数据写到数据库里，OrderItem+OrderProduct表
		5、删除购物车内已买到的商品
	*/
	var productIds []int32
	var cartList []model.ShopCart
	//                     产品id   产品数量
	productNumMap := make(map[int32]int32)
	checked := true
	r := internal.DB.Where(&model.ShopCart{AccountId: ol.AccountId, Checked: &checked}).Find(&cartList)
	if r.RowsAffected == 0 {
		//return nil, errors.New(custom_error.OrderProductList)
		return
	}
	for _, cart := range cartList {
		productIds = append(productIds, cart.ProductId)
		productNumMap[cart.ProductId] = cart.Num
	}
	productRes, err := internal.ProductClient.BatchGetProduct(context.Background(), &pb.BatchProductIdReq{Ids: productIds})
	if err != nil {
		return nil, errors.New(custom_error.ProductNotFound)
	}
	var amount float32 // 总价 = 单价*数量
	var orderProductList []model.OrderProduct
	var stockItemList []*pb.ProductStockItem

	for _, p := range productRes.ItemList {
		amount += p.RealPrice * float32(productNumMap[p.Id])
		var orderProduct = model.OrderProduct{
			ProductId:   p.Id,
			ProductName: p.Name,
			CoverImage:  p.CoverImages,
			RealPrice:   p.RealPrice,
			Num:         productNumMap[p.Id],
		}
		orderProductList = append(orderProductList, orderProduct)
		stockItem := &pb.ProductStockItem{
			ProductID: p.Id,
			Num:       productNumMap[p.Id],
		}
		stockItemList = append(stockItemList, stockItem)
	}
	_, err = internal.StockClient.Sell(context.Background(), &pb.SellItem{StockItemList: stockItemList})
	if err != nil {
		return nil, errors.New(custom_error.StockNotEnough)
	}

	tx := internal.DB.Begin()
	orderItem := model.OrderItem{
		AccountId:      ol.AccountId,
		OrderNo:        uuid.NewV4().String(),
		Status:         "unPay",
		Addr:           ol.Addr,
		Receiver:       ol.Receiver,
		ReceiverMobile: ol.Mobile,
		PostCode:       ol.PostCode,
		OrderAmount:    amount,
	}
	result := tx.Save(&orderItem)
	if result.Error != nil || result.RowsAffected < 1 {
		tx.Rollback()
		_, err = internal.StockClient.BackStock(context.Background(), &pb.SellItem{StockItemList: stockItemList})
		if err != nil {
			return nil, errors.New(custom_error.StockBackFiled)
		}
		return nil, errors.New(custom_error.CreateOrderFailed + "保存orderItem")
	}
	for i := 0; i < len(orderProductList); i++ {
		orderProductList[i].OrderId = orderItem.OrderNo
	}
	fmt.Println(orderProductList)
	result = tx.CreateInBatches(orderProductList, 50)
	if result.Error != nil || result.RowsAffected < 1 {
		tx.Rollback()
		_, err = internal.StockClient.BackStock(context.Background(), &pb.SellItem{StockItemList: stockItemList})
		if err != nil {
			return nil, errors.New(custom_error.StockBackFiled)
		}
		return nil, errors.New(custom_error.CreateOrderFailed + "赋值商品订单号")
	}
	result = tx.Where(&model.ShopCart{Checked: &checked, AccountId: item.AccountId}).Delete(&model.ShopCart{})
	if result.Error != nil || result.RowsAffected < 1 {
		tx.Rollback()
		_, err = internal.StockClient.BackStock(context.Background(), &pb.SellItem{StockItemList: stockItemList})
		if err != nil {
			return nil, errors.New(custom_error.StockBackFiled)
		}
		return nil, errors.New(custom_error.CreateOrderFailed + "更新购物车是否选中")
	}

	tx.Commit()

	return primitive.CommitMessageState
}

func (ol *OrderListener) CheckLocalTransaction(*primitive.MessageExt) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

func (s CartOrderServer) CreateOrder(ctx context.Context, item *pb.OrderItemReq) (*pb.OrderItemRes, error) {

	res := pb.OrderItemRes{
		Id:       orderItem.ID,
		OrderNum: orderItem.OrderNo,
		Amount:   orderItem.OrderAmount,
	}
	return &res, nil

}

func (s CartOrderServer) OrderList(ctx context.Context, item *pb.OrderPagingReq) (*pb.OrderListRes, error) {
	var orderList []model.OrderItem
	var res pb.OrderListRes
	var total int64
	internal.DB.Where(&model.OrderItem{
		AccountId: item.AccountId,
	}).Count(&total)
	res.Total = int32(total)

	internal.DB.Where(&model.OrderItem{
		AccountId: item.AccountId,
	}).Scopes(internal.MyPaging(int(item.PageNo), int(item.PageSize))).Find(&orderList)
	for _, item := range orderList {
		r := ConventOrderModel2Pb(item)
		res.ItemList = append(res.ItemList, r)
	}
	return &res, nil

}

func (s CartOrderServer) OrderDetail(ctx context.Context, item *pb.OrderItemReq) (*pb.OrderItemDetailRes, error) {
	var orderDetail model.OrderItem
	var detailRes pb.OrderItemDetailRes
	r := internal.DB.Where(&model.OrderItem{BaseMode: model.BaseMode{
		ID: item.Id,
	}, AccountId: item.AccountId}).First(&orderDetail)
	if r.RowsAffected == 0 {
		return nil, errors.New(custom_error.OrderNotFound)
	}
	res := ConventOrderModel2Pb(orderDetail)
	detailRes.Order = res
	var orderProductList []model.OrderProduct
	internal.DB.Where(&model.OrderProduct{OrderId: orderDetail.OrderNo}).Find(&orderProductList)
	for _, product := range orderProductList {
		detailRes.ProductList = append(detailRes.ProductList, ConvertOrderProductModel2Pb(product))
	}
	return &detailRes, nil
}

func (s CartOrderServer) ChangeOrderStatus(ctx context.Context, status *pb.OrderStatus) (*emptypb.Empty, error) {
	r := internal.DB.Model(&model.OrderItem{}).
		Where("order_no = ?", status.OrderNo).
		Update("status = ?", status.Status)
	//  update 零值问题
	if r.RowsAffected == 0 {
		return nil, errors.New(custom_error.OrderNotFound)
	}
	return &emptypb.Empty{}, nil
}

func ConventOrderModel2Pb(o model.OrderItem) *pb.OrderItemRes {
	res := pb.OrderItemRes{
		Id:        o.ID,
		AccountId: o.AccountId,
		PayType:   o.PayType,
		OrderNum:  o.OrderNo,
		PostCode:  o.PostCode,
		Amount:    o.OrderAmount,
		Addr:      o.Addr,
		Receiver:  o.Receiver,
		Mobile:    o.ReceiverMobile,
		Status:    o.Status,
	}
	return &res
}

func ConvertOrderProductModel2Pb(p model.OrderProduct) *pb.OrderProductRes {
	return &pb.OrderProductRes{
		Id:          p.ID,
		OrderId:     p.OrderId,
		ProductId:   p.ProductId,
		Num:         p.Num,
		ProductName: p.ProductName,
		RealPrice:   p.RealPrice,
		CoverImage:  p.CoverImage,
	}

}
