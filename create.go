package alipay

import (
	"encoding/json"
)

// https://docs.open.alipay.com/api_1/alipay.trade.create/

type CreateParamGoods struct {
	GoodsId       string `json:"goods_id"`       // 商品的编号
	GoodsName     string `json:"goods_name"`     // 商品名称
	Quantity      int    `json:"quantity"`       // 商品数量
	Price         string `json:"price"`          // 商品单价
	GoodsCategory string `json:"goods_category"` // 商品类目
	Body          string `json:"body"`           // 商品描述信息
	ShowUrl       string `json:"show_url"`       // 商品的展示地址
}

type CreateParamExtendParams struct {
	SysServiceProviderId string `json:"sys_service_provider_id"` // 系统商编号
}

type CreateParam struct {
	OutTradeNo         string                  `json:"out_trade_no"`        // 商户订单号
	SellerId           string                  `json:"seller_id"`           // 卖家支付宝用户ID
	TotalAmount        string                  `json:"total_amount"`        // 订单总金额
	DiscountableAmount string                  `json:"discountable_amount"` // 参与优惠计算的金额
	Subject            string                  `json:"subject"`             // 订单标题
	Body               string                  `json:"body"`                // 订单描述
	BuyerId            string                  `json:"buyer_id"`            // 买家的支付宝用户id
	GoodsDetailList    []*CreateParamGoods     `json:" goods_detail"`       // 订单包含的商品列表信息
	OperatorId         string                  `json:"operator_id"`         // 商户操作员编号
	StoreId            string                  `json:"store_id"`            // 商户门店编号
	TerminalId         string                  `json:"terminal_id"`         // 商户机具终端编号
	ExtendParams       CreateParamExtendParams `json:"extend_params"`       // 业务扩展参数
	TimeoutExpress     string                  `json:"timeout_express"`     // 该笔订单允许的最晚付款时间
	BusinessParams     string                  `json:"business_params"`     // 商户传入业务信息
}

type CreateResponse struct {
	ResponseError
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	TradeNo    string `json:"trade_no"`     // 支付宝交易号
}

func (alipay *Alipay) Create(param *CreateParam) (int, *CreateResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		MethodAlipayTradeCreate,
	)
	if err != nil {
		return 0, nil, err
	}
	createResponse := new(CreateResponse)
	if err := json.Unmarshal(body, createResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, createResponse, nil
}
