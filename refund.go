package alipay

import (
	"encoding/json"
)

type RefundParamGoodsDetail struct {
	// 商品的编号
	GoodsId string `json:" goods_id,omitempty"`
	// 支付宝定义的统一商品编号
	AlipayGoodsId string `json:"alipay_goods_id,omitempty"`
	// 商品名称
	GoodsName string `json:" goods_name,omitempty"`
	// 商品数量
	Quantity string `json:"quantity,omitempty"`
	// 商品单价，单位为元
	Price string `json:"price,omitempty"`
	// 商品类目
	GoodsCategory string `json:"goods_category,omitempty"`
	// 商品描述信息
	Body string `json:"body,omitempty"`
	// 商品的展示地址
	ShowUrl string `json:"show_url,omitempty"`
}

type RefundParam struct {
	OutTradeNo     string                    `json:"out_trade_no,omitempty"`
	TradeNo        string                    `json:"trade_no,omitempty"`
	RefundAmount   string                    `json:"refund_amount,omitempty"`
	RefundReason   string                    `json:"refund_reason,omitempty"`
	OutRequestNo   string                    `json:"out_request_no,omitempty"`
	OperatorId     string                    `json:"operator_id,omitempty"`
	StoreId        string                    `json:"store_id,omitempty"`
	TerminalId     string                    `json:"terminal_id,omitempty"`
	GoodDetailList []*RefundParamGoodsDetail `json:"goods_detail,omitempty"`
}

type RefundResponseRefundDetailItem struct {
	// 交易使用的资金渠道
	RefundChannel string `json:"refund_channel"`
	// 银行卡支付时的银行代码
	BankCode string `json:"bank_code"`
	// 该支付工具类型所使用的金额
	Amount string `json:"amount"`
	// 渠道实际付款金额
	RealAmount string `json:"real_amount"`
	// 渠道所使用的资金类型,目前只在资金渠道(fund_channel)是银行卡渠道(BANKCARD)的情况下才返回该信息(DEBIT_CARD:借记卡,CREDIT_CARD:信用卡,MIXED_CARD:借贷合一卡)
	FundType string `json:"fund_type"`
}

type RefundResponse struct {
	ResponseError
	// 支付宝交易号
	TradeNo string `json:"trade_no"`
	// 商户订单号
	OutTradeNo string `json:"out_trade_no"`
	// 用户的登录id
	BuyerLogonId string `json:"buyer_logon_id"`
	// 本次退款是否发生了资金变化
	FundChange string `json:"fund_change"`
	// 退款总金额
	RefundFee string `json:"refund_fee"`
	// 退款支付时间
	GmtRefundPay string `json:"gmt_refund_pay"`
	// 退款使用的资金渠道
	RefundDetailItemList []*RefundResponseRefundDetailItem `json:"refund_detail_item_list"`
	// 交易在支付时候的门店名称
	StoreName string `json:"store_name"`
	// 买家在支付宝的用户id
	BuyerUserId string `json:"buyer_user_id"`
	// 本次退款金额中买家退款金额
	PresentRefundBuyerAmount string `json:"present_refund_buyer_amount"`
	// 本次退款金额中平台优惠退款金额
	PresentRefundDiscountAmount string `json:"present_refund_discount_amount"`
	// 本次退款金额中商家优惠退款金额
	PresentRefundMdiscountAmount string `json:"present_refund_mdiscount_amount"`
}

func (resp *RefundResponse) IsNotEnoughBalance() bool {
	return resp.SubCode == "ACQ.SELLER_BALANCE_NOT_ENOUGH"
}

func (resp *RefundResponse) IsTradeStatusError() bool {
	return resp.SubCode == "ACQ.TRADE_STATUS_ERROR"
}

func (resp *RefundResponse) IsNotEqualTotal() bool {
	return resp.SubCode == "ACQ.REFUND_AMT_NOT_EQUAL_TOTAL"
}

func (alipay *Alipay) Refund(param *RefundParam) (int, *RefundResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_TRADE_REFUND,
	)
	if err != nil {
		return 0, nil, err
	}
	refundResponse := new(RefundResponse)
	if err := json.Unmarshal(body, refundResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, refundResponse, nil
}
