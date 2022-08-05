package model

import "time"

type OrderItem struct {
	BaseMode
	AccountId      int32  `gorm:"type:int;index"`
	OrderNo        string `gorm:"type:varchar(64);index"`
	PayType        string `gorm:"type:varchar(64)"` // int wx 银联 支付宝
	Status         string `gorm:"type:varchar(16)"` // 未支付 支付成功 超市关闭
	TradeNo        string `gorm:"type:varchar(64)"` // 保证双方对账， 第三方对账平台的凭证
	Addr           string `gorm:"type:varchar(64)"` //
	Receiver       string `gorm:"type:varchar(16)"`
	ReceiverMobile string `gorm:"type:varchar(11)"`
	PostCode       string `gorm:"type:varchar(16)"`
	OrderAmount    float32
	PayTime        *time.Time `gorm:"type:datetime"`
}

type OrderProduct struct {
	BaseMode
	OrderId     string `gorm:"type:type:varchar(64);index"`
	ProductId   int32  `gorm:"type:int;index"`
	ProductName string `gorm:"type:varchar(64);index"`
	CoverImage  string `gorm:"type:varchar(128)"`
	RealPrice   float32
	Num         int32 `gorm:"type:int"`
}
