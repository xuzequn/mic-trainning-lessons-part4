package req

import "time"

type ShopCartReq struct {
	Id        int32 `json:"id"`
	AccountId int32 `json:"accountId" binding:"required"`
	ProductId int32 `json:"productId" binding:"required"`
	Num       int32 `json:"num" binding:"required"`
	Checked   *bool `json:"checked" binding:"required"`
}

type OrderReq struct {
	Id             int32  `json:"id"`
	AccountId      int32  `json:"accountId" binding:"required"`
	OrderNo        string `json:"orderNo" binding:"required"`
	PayType        string `json:"payType"` // int wx 银联 支付宝
	Status         string `json:"status"`  // 未支付 支付成功 超市关闭
	TradeNo        string `json:"tradeNo"` // 保证双方对账， 第三方对账平台的凭证
	Addr           string `json:"addr"`
	Receiver       string `json:"receiver"`
	ReceiverMobile string `json:"receiverMobile"`
	PostCode       string `json:"postCode"`
	OrderAmount    float32
	PayTime        *time.Time `gorm:"type:datetime"`
}
