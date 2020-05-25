package alipay

import (
	"encoding/json"
)

// RefundParamGoodsDetail ...
type RefundParamGoodsDetail struct {
	GoodsID       string `json:" goods_id,omitempty"`       // 商品的编号
	AlipayGoodsID string `json:"alipay_goods_id,omitempty"` // 支付宝定义的统一商品编号
	GoodsName     string `json:" goods_name,omitempty"`     // 商品名称
	Quantity      string `json:"quantity,omitempty"`        // 商品数量
	Price         string `json:"price,omitempty"`           // 商品单价，单位为元
	GoodsCategory string `json:"goods_category,omitempty"`  // 商品类目
	Body          string `json:"body,omitempty"`            // 商品描述信息
	ShowURL       string `json:"show_url,omitempty"`        // 商品的展示地址
}

// RefundParam ...
type RefundParam struct {
	OutTradeNo     string                    `json:"out_trade_no,omitempty"`
	TradeNo        string                    `json:"trade_no,omitempty"`
	RefundAmount   string                    `json:"refund_amount,omitempty"`
	RefundReason   string                    `json:"refund_reason,omitempty"`
	OutRequestNo   string                    `json:"out_request_no,omitempty"`
	OperatorID     string                    `json:"operator_id,omitempty"`
	StoreID        string                    `json:"store_id,omitempty"`
	TerminalID     string                    `json:"terminal_id,omitempty"`
	GoodDetailList []*RefundParamGoodsDetail `json:"goods_detail,omitempty"`
}

// RefundResponseRefundDetailItem ...
type RefundResponseRefundDetailItem struct {
	RefundChannel string `json:"refund_channel"` // 交易使用的资金渠道
	BankCode      string `json:"bank_code"`      // 银行卡支付时的银行代码
	Amount        string `json:"amount"`         // 该支付工具类型所使用的金额
	RealAmount    string `json:"real_amount"`    // 渠道实际付款金额
	FundType      string `json:"fund_type"`      // 渠道所使用的资金类型,目前只在资金渠道(fund_channel)是银行卡渠道(BANKCARD)的情况下才返回该信息(DEBIT_CARD:借记卡,CREDIT_CARD:信用卡,MIXED_CARD:借贷合一卡)
}

// RefundResponse ...
type RefundResponse struct {
	ResponseError
	TradeNo                      string                            `json:"trade_no"`                        // 支付宝交易号
	OutTradeNo                   string                            `json:"out_trade_no"`                    // 商户订单号
	BuyerLogonID                 string                            `json:"buyer_logon_id"`                  // 用户的登录id
	FundChange                   string                            `json:"fund_change"`                     // 本次退款是否发生了资金变化
	RefundFee                    string                            `json:"refund_fee"`                      // 退款总金额
	GmtRefundPay                 string                            `json:"gmt_refund_pay"`                  // 退款支付时间
	RefundDetailItemList         []*RefundResponseRefundDetailItem `json:"refund_detail_item_list"`         // 退款使用的资金渠道
	StoreName                    string                            `json:"store_name"`                      // 交易在支付时候的门店名称
	BuyerUserID                  string                            `json:"buyer_user_id"`                   // 买家在支付宝的用户id
	PresentRefundBuyerAmount     string                            `json:"present_refund_buyer_amount"`     // 本次退款金额中买家退款金额
	PresentRefundDiscountAmount  string                            `json:"present_refund_discount_amount"`  // 本次退款金额中平台优惠退款金额
	PresentRefundMdiscountAmount string                            `json:"present_refund_mdiscount_amount"` // 本次退款金额中商家优惠退款金额
}

// IsNotEnoughBalance ...
func (resp *RefundResponse) IsNotEnoughBalance() bool {
	return resp.SubCode == "ACQ.SELLER_BALANCE_NOT_ENOUGH"
}

// IsTradeStatusError ...
func (resp *RefundResponse) IsTradeStatusError() bool {
	return resp.SubCode == "ACQ.TRADE_STATUS_ERROR"
}

// IsNotEqualTotal ...
func (resp *RefundResponse) IsNotEqualTotal() bool {
	return resp.SubCode == "ACQ.REFUND_AMT_NOT_EQUAL_TOTAL"
}

// Refund ...
func (alipay *Alipay) Refund(param *RefundParam) (int, *RefundResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		MethodAlipayTradeRefund,
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
