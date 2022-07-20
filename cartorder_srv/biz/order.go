package biz

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/utils"
	"mic-trainning-lesson-part4/cartorder_srv/model"
	"mic-trainning-lesson-part4/internal"
	"mic-trainning-lesson-part4/proto/pb"
)

func (s CartOrderServer) CreateOrder(ctx context.Context, item *pb.OrderItemReq) (*pb.OrderItemRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s CartOrderServer) OrderList(ctx context.Context, item *pb.OrderPagingReq) (*pb.OrderListRes, error) {
	var orderList []model.OrderItem
	var res pb.OrderListRes
	var total int64
	internal.DB.Where(&model.OrderItem{
		AccountId: item.AccountId,
	}).Count(&total)
	res.Total = int32(total)

	internal.DB.Scopes(utils.Paginate)

}

func (s CartOrderServer) OrderDetail(ctx context.Context, item *pb.OrderItemReq) (*pb.OrderItemDetailRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s CartOrderServer) ChangeOrderStatus(ctx context.Context, status *pb.OrderStatus) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
