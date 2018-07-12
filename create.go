package alipay

import (
	"encoding/json"
)

// https://docs.open.alipay.com/api_1/alipay.trade.create/

type CreateParamGoods struct {
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

type CreateParamExtendParams struct {
	// 系统商编号
	SysServiceProviderId string `json:"sys_service_provider_id"`
}

type CreateParam struct {
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
	// 订单描述
	Body string `json:"body"`
	// 买家的支付宝用户id
	BuyerId string `json:"buyer_id"`
	// 订单包含的商品列表信息
	GoodsDetailList []*CreateParamGoods `json:" goods_detail"`
	// 商户操作员编号
	OperatorId string `json:"operator_id"`
	// 商户门店编号
	StoreId string `json:"store_id"`
	// 商户机具终端编号
	TerminalId string `json:"terminal_id"`
	// 业务扩展参数
	ExtendParams CreateParamExtendParams `json:"extend_params"`
	// 该笔订单允许的最晚付款时间
	TimeoutExpress string `json:"timeout_express"`
	// 商户传入业务信息
	BusinessParams string `json:"business_params"`
}

type CreateResponse struct {
	ResponseError
	// 商户订单号
	OutTradeNo string `json:"out_trade_no"`
	// 支付宝交易号
	TradeNo string `json:"trade_no"`
}

func (alipay *Alipay) Create(param *CreateParam) (int, *CreateResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_TRADE_CREATE,
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
