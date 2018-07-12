package alipay

import (
	"encoding/json"
)

// https://docs.open.alipay.com/api_1/alipay.trade.precreate/

type PrecreateParamGoods struct {
	// 商品的编号
	GoodsId string `json:"goods_id"`
	// 商品名称
	GoodsName string `json:"goods_name"`
	// 商品数量
	Quantity int `json:"quantity"`
	// 商品单价
	Price string `json:"price"`
	// 商品类目
	GoodsCategory string `json:"goods_category"`
	// 商品描述信息
	Body string `json:"body"`
	// 商品的展示地址
	ShowUrl string `json:"show_url"`
}

type PrecreateParamExtendParams struct {
	// 系统商编号
	SysServiceProviderId string `json:"sys_service_provider_id"`
}

type PrecreateParam struct {
	// 商户订单号
	OutTradeNo string `json:"out_trade_no"`
	// 卖家支付宝用户ID
	SellerId string `json:"seller_id"`
	// 订单总金额
	TotalAmount string `json:"total_amount"`
	// 参与优惠计算的金额
	DiscountableAmount string `json:"discountable_amount"`
	// 订单标题
	Subject string `json:"subject"`
	// 订单包含的商品列表信息
	GoodsDetailList []*PrecreateParamGoods `json:" goods_detail"`
	// 订单描述
	Body string `json:"body"`
	// 商户操作员编号
	OperatorId string `json:"operator_id"`
	// 商户门店编号
	StoreId string `json:"store_id"`
	// 禁用渠道，用户不可用指定渠道支付
	DisablePayChannels string `json:"disable_pay_channels"`
	// 可用渠道，用户只能在指定渠道范围内支付
	EnablePayChannels string `json:"enable_pay_channels"`
	// 商户机具终端编号
	TerminalId string `json:"terminal_id"`
	// 业务扩展参数
	ExtendParams PrecreateParamExtendParams `json:"extend_params"`
	// 该笔订单允许的最晚付款时间
	TimeoutExpress string `json:"timeout_express"`
	// 商户传入业务信息
	BusinessParams string `json:"business_params"`
}

type PrecreateResponse struct {
	ResponseError
	// 商户订单号
	OutTradeNo string `json:"out_trade_no"`
	// 支付宝交易号
	TradeNo string `json:"trade_no"`
}

func (alipay *Alipay) Precreate(param *PrecreateParam) (int, *PrecreateResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_TRADE_PRECREATE,
	)
	if err != nil {
		return 0, nil, err
	}
	precreateResponse := new(PrecreateResponse)
	if err := json.Unmarshal(body, precreateResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, precreateResponse, nil
}
