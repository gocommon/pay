package wxpay

import (
	"encoding/json"
	"encoding/xml"
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
	NotifyURL    string
	IsProduction bool
}

// Wxpay Wxpay
type Wxpay struct {
	opt    Options
	client *wxpay.Client
}

// New New
func New(opt Options) *Wxpay {
	cli := wxpay.New(opt.AppID, opt.APIKey, opt.MchID, opt.IsProduction)

	p := &Wxpay{
		client: cli,
		opt:    opt,
	}

	return p
}

// Verify 支付回调验证签名,成功返回回调参数
func (p *Wxpay) Verify(in url.Values) (*pay.NoticeParams, error) {
	ok, err := wxpay.VerifyResponseValues(in, p.opt.APIKey)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, pay.ErrVerify
	}

	return NoticeParams(in), nil

}

// Success 回调成功返回数据
func (p *Wxpay) Success() string {

	var v = url.Values{}
	v.Set("return_code", "SUCCESS")
	v.Set("return_msg", "OK")

	return wxpay.URLValueToXML(v)
}

// Call 调起支付用到的数据
// form -> 自动提交form表单 html
// app -> 调起app用到的url参数
// qrcode -> 二维码图片地址
// h5 -> 自动提交form表单 html
func (p *Wxpay) Call(way pay.Way, in pay.Order) (string, error) {
	switch way {
	case pay.WayQrcode:
		// 二维码
		return p.qrcodeCall(in)
	case pay.WayApp:
		// app
		return p.appCall(in)
	case pay.WayJSAPI:
		// 公众号
		return p.jsAPICall(in)
	case pay.WayWap:
		// 手机浏览器
		return p.wapCall(in)
	case pay.WayWXXCX:
		// 小程序
		return p.wxxcxCall(in)

	}

	return "", pay.ErrWayNotDefine
}

// qrcodeCall 返回二维码地址 ip 传服务器端ip
func (p *Wxpay) qrcodeCall(in pay.Order) (string, error) {
	info, err := p.client.NativePay(wxpay.UnifiedOrderParam{
		NotifyURL:      p.opt.NotifyURL,           // 是 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
		Body:           in.Title,                  // 是 商品简单描述，该字段请按照规范传递，具体请见参数规定
		OutTradeNo:     in.ID,                     // 是 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。详见商户订单号
		TotalFee:       int(in.Amount),            // 是 订单总金额，单位为分，详见支付金额
		SpbillCreateIP: in.IP,                     // 是 APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
		TradeType:      wxpay.K_TRADE_TYPE_NATIVE, // 是 取值如下：JSAPI，NATIVE，APP等，说明详见参数规定
		ProductId:      in.ID,                     // 否 trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	})
	if err != nil {
		return "", err
	}

	return info.CodeURL, nil
}

// appCall 返回app调起支付的参数
func (p *Wxpay) appCall(in pay.Order) (string, error) {
	info, err := p.client.AppPay(p.opt.AppID, wxpay.UnifiedOrderParam{
		NotifyURL:      p.opt.NotifyURL,        // 是 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
		Body:           in.Title,               // 是 商品简单描述，该字段请按照规范传递，具体请见参数规定
		OutTradeNo:     in.ID,                  // 是 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。详见商户订单号
		TotalFee:       int(in.Amount),         // 是 订单总金额，单位为分，详见支付金额
		SpbillCreateIP: in.IP,                  // 是 APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
		TradeType:      wxpay.K_TRADE_TYPE_APP, // 是 取值如下：JSAPI，NATIVE，APP等，说明详见参数规定
		ProductId:      in.ID,                  // 否 trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	})
	if err != nil {
		return "", err
	}

	return info.Params, nil
}

// jsAPICall 返回跳转的url地址
func (p *Wxpay) jsAPICall(in pay.Order) (string, error) {
	resp, err := p.client.JSAPIPay(p.opt.AppID, wxpay.UnifiedOrderParam{
		NotifyURL:      p.opt.NotifyURL,          // 是 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
		Body:           in.Title,                 // 是 商品简单描述，该字段请按照规范传递，具体请见参数规定
		OutTradeNo:     in.ID,                    // 是 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。详见商户订单号
		TotalFee:       int(in.Amount),           // 是 订单总金额，单位为分，详见支付金额
		SpbillCreateIP: in.IP,                    // 是 APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
		TradeType:      wxpay.K_TRADE_TYPE_JSAPI, // 是 取值如下：JSAPI，NATIVE，APP等，说明详见参数规定
		OpenId:         in.OpenID,
	})
	if err != nil {
		return "", err
	}

	return resp.Params, nil
}

// wapCall 返回跳转的url地址
func (p *Wxpay) wapCall(in pay.Order) (string, error) {
	sInfo := H5SceneInfo{
		H5Info: H5Info{
			Type:    "Wap",
			WapURL:  "//",
			WapName: in.Title,
		},
	}

	d, _ := json.Marshal(sInfo)

	resp, err := p.client.WebPay(wxpay.UnifiedOrderParam{
		Body:           in.Title,                // 是 商品简单描述，该字段请按照规范传递，具体请见参数规定
		OutTradeNo:     in.ID,                   // 是 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。详见商户订单号
		TotalFee:       int(in.Amount),          // 是 订单总金额，单位为分，详见支付金额
		SpbillCreateIP: "",                      // 是 APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
		TradeType:      wxpay.K_TRADE_TYPE_MWEB, // 是 取值如下：JSAPI，NATIVE，APP等，说明详见参数规定
		NotifyURL:      p.opt.NotifyURL,         // 是 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
		SceneInfo:      string(d),
	})
	if err != nil {
		return "", err
	}

	return resp.MWebURL, nil
}

// wxxcxCall 返回跳转的url地址
func (p *Wxpay) wxxcxCall(in pay.Order) (string, error) {
	resp, err := p.client.MiniAppPay(p.opt.AppID, wxpay.UnifiedOrderParam{
		NotifyURL:      p.opt.NotifyURL,          // 是 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
		Body:           in.Title,                 // 是 商品简单描述，该字段请按照规范传递，具体请见参数规定
		OutTradeNo:     in.ID,                    // 是 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。详见商户订单号
		TotalFee:       int(in.Amount),           // 是 订单总金额，单位为分，详见支付金额
		SpbillCreateIP: in.IP,                    // 是 APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
		TradeType:      wxpay.K_TRADE_TYPE_JSAPI, // 是 取值如下：JSAPI，NATIVE，APP等，说明详见参数规定
		OpenId:         in.OpenID,
	})
	if err != nil {
		return "", err
	}

	return resp.Params, nil
}

// H5SceneInfo H5SceneInfo
type H5SceneInfo struct {
	H5Info H5Info `json:"h5_info"`
}

// H5Info H5Info
type H5Info struct {
	Type    string `json:"type"`
	WapURL  string `json:"wap_url"`
	WapName string `json:"wap_name"`
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

// BodyToValues 转request.Body的xml内容到url.Values
func BodyToValues(body string) (url.Values, error) {
	var param = make(wxpay.XMLMap)
	err := xml.Unmarshal([]byte(body), &param)
	if err != nil {
		return nil, err
	}

	return url.Values(param), nil
}
