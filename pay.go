package pay

import (
	"errors"
	"net/url"
)

// Way 支付方式
type Way string

const (
	// WayForm form支付方式，用于PC网页表单跳转到支付宝支付
	WayForm Way = "form"
	// WayQrcode 二维码支付方式，用于网页显示二维，用户扫码支付
	WayQrcode Way = "qrcode"
	// WayApp 调起支付宝app支付
	WayApp Way = "app"
	// WayWap 手机浏览器支付
	WayWap Way = "wap"
)

var (
	// ErrWayNotDefine ErrWayNotDefine
	ErrWayNotDefine = errors.New("payway not define")
)

// Payer Payer
type Payer interface {
	// Verify 支付回调验证签名,成功返回回调参数
	Verify(url.Values) (*NoticeParams, error)
	// 支付状态

	// Success 回调成功返回数据
	Success() string

	// Call 调起支付用到的数据
	// form -> 自动提交form表单 html
	// app -> 调起app用到的url参数
	// qrcode -> 二维码图片地址
	// h5 -> 自动提交form表单 html
	Call(Way, Order) (string, error)
}

// NoticeParams 回调参数
type NoticeParams struct {
	OrderID     string      // 商品订单
	PaymentID   string      // 支付单号
	TradeStatus TradeStatus //支付状态
	Amount      int32       // 支付金额
}

// Order 订单信息
type Order struct {
	ID     string // 订单ID
	Title  string // 订单详情
	Amount int32  // 支付金额 单位分
}

// TradeType 支付方式
type TradeType string

const (
	// TradeTypeApp app支付
	TradeTypeApp = "app"
	// TradeTypeForm 表单支付
	TradeTypeForm = "form"
	// TradeTypeH5 h5支付
	TradeTypeH5 = "h5"
	// TradeTypeQrcode 扫码支付
	TradeTypeQrcode = "qrcode"
)

// TradeStatus 交易状态
type TradeStatus int

const (
	// TradeStatusWait 交易创建，等待买家付款 不通知
	TradeStatusWait TradeStatus = iota
	// TradeStatusSuccess 交易支付成功
	TradeStatusSuccess
	// TradeStatusClosed 未付款交易超时关闭，或支付完成后全额退款
	TradeStatusClosed
	// TradeStatusFinished 交易结束，不可退款
	TradeStatusFinished
)
