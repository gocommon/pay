package wxpay

import (
	"net/url"
	"strconv"

	"github.com/gocommon/pay"
	"github.com/smartwalle/wxpay"
)

var _ pay.Payer = &Wxpay{}

// Options Options
type Options struct {
	AppID        string
	APIKey       string
	MchID        string
	APIDomain    string
	NotifyURL    string
	IsProduction bool
}

// Wxpay Wxpay
type Wxpay struct {
	opt    Options
	client *wxpay.Client
}

// New New
func New(opt Options) (*Wxpay, error) {
	cli := wxpay.New(opt.AppID, opt.APIKey, opt.MchID, opt.IsProduction)

	p := &Wxpay{
		client: cli,
		opt:    opt,
	}

	return p, nil
}

// Verify 支付回调验证签名,成功返回回调参数
func (p *Wxpay) Verify(in url.Values) (*pay.NoticeParams, error) {
	return nil, nil
}

// 支付状态

// Success 回调成功返回数据
func (p *Wxpay) Success() string {
	return "success"
}

// Call 调起支付用到的数据
// form -> 自动提交form表单 html
// app -> 调起app用到的url参数
// qrcode -> 二维码图片地址
// h5 -> 自动提交form表单 html
func (p *Wxpay) Call(way pay.Way, in pay.Order) (string, error) {
	switch way {
	case pay.WayQrcode:
		return p.qrcodeCall(in)
	case pay.WayApp:
		return p.appCall(in)
	case pay.WayForm:
		return p.formCall(in)
	case pay.WayWap:
		return p.wapCall(in)

	}

	return "", pay.ErrWayNotDefine
}

// wapCall 返回跳转的url地址
func (p *Wxpay) wapCall(in pay.Order) (string, error) {
	return "", nil
}

// formCall 返回跳转的url地址
func (p *Wxpay) formCall(in pay.Order) (string, error) {

	return "", nil
}

// appCall 返回app调起支付的参数
func (p *Wxpay) appCall(in pay.Order) (string, error) {

	return "", nil
}

// qrcodeCall 返回二维码地址
func (p *Wxpay) qrcodeCall(in pay.Order) (string, error) {
	return "", nil
}

// NoticeParams NoticeParams
func NoticeParams(val url.Values) *pay.NoticeParams {
	var (
		amount = val.Get("total_amount")
		status = val.Get("trade_status")

		payStatus pay.TradeStatus
	)

	amountf, _ := strconv.ParseFloat(amount, 10)
	// if err != nil {
	// 	return nil, err
	// }

	switch status {
	case "WAIT_BUYER_PAY":
		payStatus = pay.TradeStatusWait
	case "TRADE_CLOSED":
		payStatus = pay.TradeStatusClosed
	case "TRADE_SUCCESS":
		payStatus = pay.TradeStatusSuccess
	case "TRADE_FINISHED":
		payStatus = pay.TradeStatusFinished
	}

	return &pay.NoticeParams{
		OrderID:     val.Get("out_trade_no"),
		PaymentID:   val.Get("trade_no"), // 支付单号
		TradeStatus: payStatus,           //支付状态
		Amount:      int32(amountf * 100),
	}
}
