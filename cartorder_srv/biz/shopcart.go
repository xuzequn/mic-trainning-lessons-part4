package biz

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"mic-trainning-lesson-part4/cartorder_srv/model"
	"mic-trainning-lesson-part4/custom_error"
	"mic-trainning-lesson-part4/internal"
	"mic-trainning-lesson-part4/proto/pb"
)

type CartOrderServer struct {
}

func (s CartOrderServer) ShopCartItemList(ctx context.Context, req *pb.AccountReq) (*pb.CartItemListRes, error) {
	var cartItemList []model.ShopCart
	var res *pb.CartItemListRes
	var cartItemListPb []*pb.CartItemRes
	res = new(pb.CartItemListRes)
	r := internal.DB.Where(&model.ShopCart{AccountId: req.AccountId}).Find(&cartItemList)
	if r.Error != nil {
		return nil, errors.New(custom_error.ParamError)
	}
	if r.RowsAffected < 1 {
		return res, nil
	}
	for _, item := range cartItemList {
		itemPb := ConvertShopCartModel2Pb(item)
		cartItemListPb = append(cartItemListPb, itemPb)
	}

	res.ItemList = cartItemListPb
	res.Total = int32(r.RowsAffected)

	return res, nil
}

func (s CartOrderServer) AddShopCartItem(ctx context.Context, req *pb.ShopCartReq) (*pb.CartItemRes, error) {
	/*
		if 没有productId {
					添加
					} else {
					更新数量
			}
	*/
	var cart model.ShopCart
	r := internal.DB.Where(&model.ShopCart{
		AccountId: req.AccountId,
		ProductId: req.ProductId,
	}).First(&cart)
	if r.RowsAffected < 1 {
		cart.AccountId = req.AccountId
		cart.ProductId = req.ProductId
		cart.Num = req.Num
		cart.Checked = &req.Checked
	} else {
		cart.Num += req.Num
		cart.Checked = &req.Checked
	}
	internal.DB.Save(&cart)
	res := ConvertShopCartModel2Pb(cart)
	return res, nil

}

func (s CartOrderServer) UpdateShopCartItem(ctx context.Context, req *pb.ShopCartReq) (*emptypb.Empty, error) {
	var cart model.ShopCart
	r := internal.DB.Where(&model.ShopCart{
		AccountId: req.AccountId,
		ProductId: req.ProductId,
	}).Find(&cart)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CartNotFound)
	}
	if req.Num < 1 {
		return nil, errors.New(custom_error.ParamError)
	}

	cart.Num = req.Num
	cart.Checked = &req.Checked
	fmt.Println(cart)
	r = internal.DB.Updates(cart)
	if r.RowsAffected == 0 {
		zap.S().Info("更新购物车失败")
	}
	return &emptypb.Empty{}, nil
}

func (s CartOrderServer) DeleteShopCartItem(ctx context.Context, req *pb.DelShopCartItem) (*emptypb.Empty, error) {
	var cart model.ShopCart
	r := internal.DB.Where("product_id=? and account_id=?", req.ProductId, req.AccountId).Delete(&cart)
	if r.RowsAffected < 1 {
		return nil, errors.New(custom_error.CartNotFound)
	}
	return &emptypb.Empty{}, nil
}

func ConvertShopCartModel2Pb(s model.ShopCart) *pb.CartItemRes {
	cart := &pb.CartItemRes{
		Id:        s.ID,
		AccountId: s.AccountId,
		ProductId: s.ProductId,
		Num:       s.Num,
		Checked:   *s.Checked,
	}
	return cart
}
