package alipay

import (
	"encoding/json"
)

// https://docs.open.alipay.com/api_1/alipay.trade.query

type QueryParam struct {
	OutTradeNo string `json:"out_trade_no,omitempty"`
	TradeNo    string `json:"trade_no,omitempty"`
}

type QueryResponseFundBill struct {
	// 交易使用的资金渠道
	FundChannel string `json:"fund_channel"`
	// 该支付工具类型所使用的金额
	Amount string `json:"amount"`
	// 渠道实际付款金额
	RealAmount string `json:"real_amount"`
}

type QueryResponse struct {
	ResponseError

	// 支付宝交易号
	TradeNo string `json:"trade_no"`

	// 商家订单号
	OutTradeNo string `json:"out_trade_no"`

	// 买家支付宝账号
	BuyerLogonId string `json:"buyer_logon_id"`

	// 交易状态
	TradeStatus string `json:"trade_status"`

	// 交易的订单金额
	TotalAmount string `json:"total_amount"`

	// 实收金额
	ReceiptAmount string `json:"receipt_amount"`

	// 买家实付金额
	BuyerPayAmount string `json:"buyer_pay_amount"`

	PointAmount string `json:"point_amount"`

	// 交易中用户支付的可开具发票的金额
	InvoiceAmount string `json:"invoice_amount"`

	// 本次交易打款给卖家的时间
	SendPayDate string `json:"send_pay_date"`

	// 商户门店编号
	StoreId string `json:"store_id"`

	// 商户机具终端编号
	TerminalId string `json:"terminal_id"`

	// 交易支付使用的资金渠道
	FundBillList []*QueryResponseFundBill `json:"fund_bill_list"`

	// 请求交易支付中的商户店铺的名称
	StoreName string `json:"store_name"`

	// 买家在支付宝的用户id
	BuyerUserId string `json:"buyer_user_id"`

	// 买家用户类型
	BuyerUserType string `json:"buyer_user_type"`
}

func (alipay *Alipay) Query(param *QueryParam) (int, *QueryResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_TRADE_QUERY,
	)
	if err != nil {
		return 0, nil, err
	}
	queryResponse := new(QueryResponse)
	if err := json.Unmarshal(body, queryResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, queryResponse, nil
}
