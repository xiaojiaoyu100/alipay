// Package alipay https://docs.open.alipay.com/api_1/alipay.trade.query
package alipay

import (
	"encoding/json"
)

// QueryParam ...
type QueryParam struct {
	OutTradeNo string `json:"out_trade_no,omitempty"`
	TradeNo    string `json:"trade_no,omitempty"`
}

// QueryResponseFundBill ...
type QueryResponseFundBill struct {
	FundChannel string `json:"fund_channel"` // 交易使用的资金渠道
	Amount      string `json:"amount"`       // 该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount"`  // 渠道实际付款金额
}

// QueryResponse ...
type QueryResponse struct {
	ResponseError
	TradeNo        string                   `json:"trade_no"`         // 支付宝交易号
	OutTradeNo     string                   `json:"out_trade_no"`     // 商家订单号
	BuyerLogonID   string                   `json:"buyer_logon_id"`   // 买家支付宝账号
	TradeStatus    string                   `json:"trade_status"`     // 交易状态
	TotalAmount    string                   `json:"total_amount"`     // 交易的订单金额
	ReceiptAmount  string                   `json:"receipt_amount"`   // 实收金额
	BuyerPayAmount string                   `json:"buyer_pay_amount"` // 买家实付金额
	PointAmount    string                   `json:"point_amount"`
	InvoiceAmount  string                   `json:"invoice_amount"`  // 交易中用户支付的可开具发票的金额
	SendPayDate    string                   `json:"send_pay_date"`   // 本次交易打款给卖家的时间
	StoreID        string                   `json:"store_id"`        // 商户门店编号
	TerminalID     string                   `json:"terminal_id"`     // 商户机具终端编号
	FundBillList   []*QueryResponseFundBill `json:"fund_bill_list"`  // 交易支付使用的资金渠道
	StoreName      string                   `json:"store_name"`      // 请求交易支付中的商户店铺的名称
	BuyerUserID    string                   `json:"buyer_user_id"`   // 买家在支付宝的用户id
	BuyerUserType  string                   `json:"buyer_user_type"` // 买家用户类型
}

// Query ...
func (alipay *Alipay) Query(param *QueryParam) (int, *QueryResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		MethodAlipayTradeQuery,
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
