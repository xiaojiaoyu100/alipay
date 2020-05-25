// Package alipay https://docs.open.alipay.com/api_1/alipay.trade.cancel/
package alipay

import (
	"encoding/json"
)

// CancelParam ...
type CancelParam struct {
	TradeNo    string `json:"trade_no"`     // 商户订单号
	OutTradeNo string `json:"out_trade_no"` // 支付宝交易号
}

// CancelResponse ...
type CancelResponse struct {
	ResponseError
	TradeNo    string `json:"trade_no"`     // 支付宝交易号
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	RetryFlag  string `json:"retry_flag"`   // 是否需要重试
	Action     string `json:"action"`       // 本次撤销触发的交易动作
}

// Cancel ...
func (alipay *Alipay) Cancel(param CancelParam) (int, *CancelResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		MethodAlipayTradeCancel,
	)
	if err != nil {
		return 0, nil, err
	}
	cancelResponse := new(CancelResponse)
	if err := json.Unmarshal(body, cancelResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, cancelResponse, nil
}
