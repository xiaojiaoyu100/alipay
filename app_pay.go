// https://docs.open.alipay.com/204/105465/
package alipay

import (
	"errors"
	"log"
	"net/url"
)

type AppPayParamExtUserInfo struct {
	Name          string `json:"name"`            // 姓名
	Mobile        string `json:"mobile"`          // 手机号
	CertType      string `json:"cert_type"`       // 身份证：IDENTITY_CARD、护照：PASSPORT、军官证：OFFICER_CARD、士兵证：SOLDIER_CARD、户口本：HOKOU等。如有其它类型需要支持，请与蚂蚁金服工作人员联系。
	CertNo        string `json:"cert_no"`         // 证件号
	MinAge        string `json:"min_age"`         // 允许的最小买家年龄，买家年龄必须大于等于所传数值 注：  1. need_check_info=T时该参数才有效  2. min_age为整数，必须大于等于0
	FixBuyer      string `json:"fix_buyer"`       // 是否强制校验付款人身份信息 T:强制校验，F：不强制
	NeedCheckInfo string `json:"need_check_info"` // 是否强制校验身份信息 T:强制校验，F：不强制
}

type AppPayParamExtendParams struct {
	SysServiceProviderId string `json:"sys_service_provider_id,omitempty"` // 系统商编号，该参数作为系统商返佣数据提取的依据，请填写系统商签约协议的PID
	NeedBuyerRealnamed   string `json:"needBuyerRealnamed,omitempty"`      //是否发起实名校验 T：发起 F：不发起
	TransMemo            string `json:"TRANS_MEMO,omitempty"`              // 账务备注 注：该字段显示在离线账单的账务备注中
	HbFqNum              string `json:"hb_fq_num,omitempty"`               // 花呗分期数（目前仅支持3、6、12） 注：使用该参数需要仔细阅读“花呗分期接入文档”
	HbFqSellerPercent    string `json:"hb_fq_seller_percent,omitempty"`    // 卖家承担收费比例，商家承担手续费传入100，用户承担手续费传入0，仅支持传入100、0两种，其他比例暂不支持  注：使用该参数需要仔细阅读“花呗分期接入文档”
}

type AppPayParam struct {
	Body               string                 `json:"body,omitempty"`                 // 对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body。
	Subject            string                 `json:"subject"`                        // 商品的标题/交易标题/订单标题/订单关键字等。
	OutTradeNo         string                 `json:"out_trade_no"`                   // 商户网站唯一订单号
	TimeoutExpress     string                 `json:"timeout_express,omitempty"`      // 该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。 注：若为空，则默认为15d。
	TotalAmount        string                 `json:"total_amount"`                   // 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	ProductCode        string                 `json:"product_code"`                   // 销售产品码，商家和支付宝签约的产品码，为固定值QUICK_MSECURITY_PAY
	GoodsType          string                 `json:"goods_type,omitempty"`           // 商品主类型：0—虚拟类商品，1—实物类商品 注：虚拟类商品不支持使用花呗渠道
	PassbackParams     string                 `json:"passback_params,omitempty"`      // 公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝
	PromoParams        string                 `json:"promo_params"`                   // 优惠参数 注：仅与支付宝协商后可用
	ExtendParams       string                 `json:"extend_params,omitempty"`        // 业务扩展参数，详见下面的“业务扩展参数说明”
	EnablePayChannels  string                 `json:"enable_pay_channels,omitempty"`  // 可用渠道，用户只能在指定渠道范围内支付 当有多个渠道时用“,”分隔 注：与disable_pay_channels互斥
	DisablePayChannels string                 `json:"disable_pay_channels,omitempty"` // 禁用渠道，用户不可用指定渠道支付 当有多个渠道时用“,”分隔  注：与enable_pay_channels互斥
	StoreId            string                 `json:"store_id,omitempty"`             // 商户门店编号。该参数用于请求参数中以区分各门店，非必传项。
	ExtUserInfo        AppPayParamExtUserInfo `json:"ext_user_info,omitempty"`        // 外部指定买家，详见外部用户ExtUserInfo参数说明
}

func (alipay *Alipay) AppPay(param *AppPayParam, notifyUrl string) (string, error) {
	param.ProductCode = "QUICK_MSECURITY_PAY"

	if len(param.OutTradeNo) == 0 {
		text := "商户订单号不能为空"
		log.Println(text)
		return "", errors.New(text)
	}

	if len(param.TotalAmount) == 0 {
		text := "订单金额不能为空"
		log.Println(text)
		return "", errors.New(text)
	}

	if len(param.Subject) == 0 {
		text := "订单标题不能为空"
		log.Println(text)
		return "", errors.New(text)
	}

	if len(param.TimeoutExpress) == 0 {
		text := "订单允许的最晚付款时间不能为空"
		log.Println(text)
		return "", errors.New(text)
	}

	paramStr, err := alipay.MakeParam(
		param,
		MethodAlipayTradeAppPay,
		WithNotifyUrl(notifyUrl),
	)
	if err != nil {
		log.Println("支付宝app支付构造参数失败: ", err)
		return "", err
	}
	url, err := url.Parse(AlipayGateway)
	if err != nil {
		log.Println("解析支付宝网关失败: ", err)
		return "", err
	}
	url.RawQuery = paramStr
	return url.String(), nil
}
