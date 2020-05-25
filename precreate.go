// Package alipay https://docs.open.alipay.com/api_1/alipay.trade.precreate/
package alipay

import (
	"encoding/json"
)

// PrecreateParamGoods ...
type PrecreateParamGoods struct {
	GoodsID       string `json:"goods_id"`       // 商品的编号
	GoodsName     string `json:"goods_name"`     // 商品名称
	Quantity      int    `json:"quantity"`       // 商品数量
	Price         string `json:"price"`          // 商品单价
	GoodsCategory string `json:"goods_category"` // 商品类目
	Body          string `json:"body"`           // 商品描述信息
	ShowURL       string `json:"show_url"`       // 商品的展示地址
}

// PrecreateParamExtendParams ...
type PrecreateParamExtendParams struct {
	SysServiceProviderID string `json:"sys_service_provider_id"` // 系统商编号
}

// PrecreateParam ...
type PrecreateParam struct {
	OutTradeNo         string                     `json:"out_trade_no"`         // 商户订单号
	SellerID           string                     `json:"seller_id"`            // 卖家支付宝用户ID
	TotalAmount        string                     `json:"total_amount"`         // 订单总金额
	DiscountableAmount string                     `json:"discountable_amount"`  // 参与优惠计算的金额
	Subject            string                     `json:"subject"`              // 订单标题
	GoodsDetailList    []*PrecreateParamGoods     `json:" goods_detail"`        // 订单包含的商品列表信息
	Body               string                     `json:"body"`                 // 订单描述
	OperatorID         string                     `json:"operator_id"`          // 商户操作员编号
	StoreID            string                     `json:"store_id"`             // 商户门店编号
	DisablePayChannels string                     `json:"disable_pay_channels"` // 禁用渠道，用户不可用指定渠道支付
	EnablePayChannels  string                     `json:"enable_pay_channels"`  // 可用渠道，用户只能在指定渠道范围内支付
	TerminalID         string                     `json:"terminal_id"`          // 商户机具终端编号
	ExtendParams       PrecreateParamExtendParams `json:"extend_params"`        // 业务扩展参数
	TimeoutExpress     string                     `json:"timeout_express"`      // 该笔订单允许的最晚付款时间
	BusinessParams     string                     `json:"business_params"`      // 商户传入业务信息
}

// PrecreateResponse ...
type PrecreateResponse struct {
	ResponseError
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	TradeNo    string `json:"trade_no"`     // 支付宝交易号
}

// Precreate ...
func (alipay *Alipay) Precreate(param *PrecreateParam) (int, *PrecreateResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		MethodAlipayTradePrecreate,
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
