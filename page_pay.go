package alipay

import (
	"log"
	"net/url"
)

// https://docs.open.alipay.com/270/alipay.trade.page.pay

type PagePayParamExtendParams struct {
	// 系统商编号
	SysServiceProviderId string `json:"sys_service_provider_id,omitempty"`
	// 花呗分期数（目前仅支持3、6、12）注：使用该参数需要仔细阅读“花呗分期接入文档”
	HbFqNum string `json:"hb_fq_num,omitempty"`
	// 卖家承担收费比例，商家承担手续费传入100，用户承担手续费传入0，仅支持传入100、0两种，其他比例暂不支持注：使用该参数需要仔细阅读“花呗分期接入文档”
	HbFqSellerPercent string `json:"hb_fq_seller_percent,omitempty"`
}

type PagePayParam struct {
	// 商户订单号，64个字符以内、可包含字母、数字、下划线；需保证在商户端不重复
	OutTradeNo string `json:"out_trade_no"`
	// 销售产品码，与支付宝签约的产品码名称。 注：目前仅支持FAST_INSTANT_TRADE_PAY
	ProductCode string `json:"product_code"`
	// 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	TotalAmount string `json:"total_amount"`
	// 订单标题
	Subject string `json:"subject"`
	// 订单描述
	Body string `json:"body,omitempty"`
	// 订单包含的商品列表信息，Json格式
	GoodsDetail string `json:"goods_detail,omitempty"`
	// 公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝只会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝
	PassbackParams string `json:"passback_params,omitempty"`
	// 业务扩展参数，详见业务扩展参数说明
	ExtendParams string `json:"extend_params,omitempty"`
	// 商品主类型：0&mdash;虚拟类商品，1&mdash;实物类商品（默认） 注：虚拟类商品不支持使用花呗渠道
	GoodsType string `json:"goods_type,omitempty"`
	// 该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。 该参数在请求到支付宝时开始计时。
	TimeoutExpress string `json:"timeout_express,omitempty"`
	// 可用渠道，用户只能在指定渠道范围内支付
	EnablePayChannels string `json:"enable_pay_channels,omitempty"`
	// 禁用渠道
	DisablePayChannels string `json:"disable_pay_channels,omitempty"`
	// 获取用户授权信息
	AuthToken string `json:"auth_token,omitempty"`
	// PC扫码支付的方式
	QrPayMode string `json:"qr_pay_mode,omitempty"`
	// 商户自定义二维码宽度
	QrCodeWidth string `json:"qrcode_width,omitempty"`
}

func (alipay *Alipay) PagePay(param *PagePayParam, notifyUrl string, returnUrl string) (string, error) {
	param.ProductCode = "FAST_INSTANT_TRADE_PAY"

	if len(param.OutTradeNo) == 0 {
		panic("商户订单号不能为空")
	}

	if len(param.TotalAmount) == 0 {
		panic("订单金额不能为空")
	}

	if len(param.Subject) == 0 {
		panic("订单标题不能为空")
	}

	if len(param.TimeoutExpress) == 0 {
		panic("订单允许的最晚付款时间不能为空")
	}

	paramStr, err := alipay.MakeParam(
		param,
		METHOD_ALIPAY_TRADE_PAGE_PAY,
		WithNotifyUrl(notifyUrl),
		WithReturnUrl(returnUrl),
	)
	if err != nil {
		log.Println("支付宝电脑网站支付构造参数失败: ", err)
		return "", err
	}
	url, err := url.Parse(ALIPAY_GATEWAY)
	if err != nil {
		log.Println("解析支付宝网关失败: ", err)
		return "", err
	}
	url.RawQuery = paramStr
	return url.String(), nil
}
