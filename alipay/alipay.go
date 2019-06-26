package alipay

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gocommon/pay"
	"github.com/smartwalle/alipay"
)

var _ pay.Payer = &Alipay{}

// Options Options
type Options struct {
	AppID         string
	AliPublicKey  string // ali公钥
	AppPrivateKey string // 应用私钥
	IsProduction  bool
	NotifyURL     string // 异步回调地址
	ReturnURL     string // 同步回调地址
}

// Alipay Alipay
type Alipay struct {
	opt    Options
	client *alipay.Client
}

// New New
func New(opt Options) (*Alipay, error) {
	cli, err := alipay.New(opt.AppID, opt.AliPublicKey, opt.AppPrivateKey, opt.IsProduction)
	if err != nil {
		return nil, err
	}

	p := &Alipay{
		client: cli,
		opt:    opt,
	}

	return p, nil
}

// Verify 支付回调验证签名,成功返回回调参数
func (p *Alipay) Verify(in url.Values) (*pay.NoticeParams, error) {
	ok, err := p.client.VerifySign(in)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("verify failed")
	}

	return NoticeParams(in), nil
}

// 支付状态

// Success 回调成功返回数据
func (p *Alipay) Success() string {
	return "success"
}

// Call 调起支付用到的数据
// form -> 自动提交form表单 html
// app -> 调起app用到的url参数
// qrcode -> 二维码图片地址
// h5 -> 自动提交form表单 html
func (p *Alipay) Call(way pay.Way, in pay.Order) (string, error) {
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
func (p *Alipay) wapCall(in pay.Order) (string, error) {
	u, err := p.client.TradeWapPay(alipay.TradeWapPay{
		Trade: alipay.Trade{
			Subject:     in.Title,
			OutTradeNo:  in.ID,
			TotalAmount: fmt.Sprintf("%.2f", float32(in.Amount)/100),
			NotifyURL:   p.opt.NotifyURL,
			ReturnURL:   p.opt.ReturnURL,
		},
	})

	if err != nil {
		return "", err
	}

	return u.String(), nil
}

// formCall 返回跳转的url地址
func (p *Alipay) formCall(in pay.Order) (string, error) {

	u, err := p.client.TradePagePay(alipay.TradePagePay{
		Trade: alipay.Trade{
			Subject:     in.Title,
			OutTradeNo:  in.ID,
			TotalAmount: fmt.Sprintf("%.2f", float32(in.Amount)/100),
			NotifyURL:   p.opt.NotifyURL,
			ReturnURL:   p.opt.ReturnURL,
		},
	})

	if err != nil {
		return "", err
	}

	return u.String(), nil
}

// appCall 返回app调起支付的参数
func (p *Alipay) appCall(in pay.Order) (string, error) {

	return p.client.TradeAppPay(alipay.TradeAppPay{
		Trade: alipay.Trade{
			Subject:     in.Title,
			OutTradeNo:  in.ID,
			TotalAmount: fmt.Sprintf("%.2f", float32(in.Amount)/100),
		},
		// TimeExpire: "1d", // 该笔订单允许的最晚付款时间，逾期将关闭交易 取值范围：1m～15d 不接受小数点
	})
}

// qrcodeCall 返回二维码地址
func (p *Alipay) qrcodeCall(in pay.Order) (string, error) {
	resp, err := p.client.TradePreCreate(alipay.TradePreCreate{
		Trade: alipay.Trade{
			Subject:     in.Title,
			OutTradeNo:  in.ID,
			TotalAmount: fmt.Sprintf("%.2f", float32(in.Amount)/100),
		},
	})
	if err != nil {
		return "", err
	}

	if !resp.IsSuccess() {
		return "", errors.New("sign failed")
	}

	return resp.Content.QRCode, nil
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
