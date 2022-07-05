package req

type ShopCartReq struct {
	Id        int32 `json:"id"`
	AccountId int32 `json:"accountId" binding:"required"`
	ProductId int32 `json:"productId" binding:"required"`
	Num       int32 `json:"num" binding:"required"`
	Checked   *bool `json:"checked" binding:"required"`
}
