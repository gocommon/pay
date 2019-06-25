package pay

// Payer Payer
type Payer interface {
	Signer
	Caller
}

// Signer Signer
type Signer interface {
	// Sign 支付回调验证签名,
	Sign(contentType string, body []byte) error
}

// Caller Caller
type Caller interface {
	// Params 调起支付用到的参数
	Params(Order) string
}

// Order Order
type Order struct {
	ID     string // 订单ID
	Info   string // 订单详情
	Amount int32  // 支付金额 单位分
}
