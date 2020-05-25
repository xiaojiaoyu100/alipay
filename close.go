package alipay

import (
	"encoding/json"
)

// CloseParam ...
type CloseParam struct {
	TradeNo    string `json:"trade_no,omitempty"`     // 支付宝交易号
	OutTradeNo string `json:"out_trade_no,omitempty"` // 商户订单号
	OperatorID string `json:"operator_id,omitempty"`  // 卖家端自定义的的操作员 ID
}

// CloseResponse ...
type CloseResponse struct {
	ResponseError
	TradeNo    string `json:"trade_no"`     // 支付宝交易号
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
}

// Close ...
func (alipay *Alipay) Close(param *CloseParam) (int, *CloseResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		MethodAlipayTradeClose,
	)
	if err != nil {
		return 0, nil, err
	}
	closeResponse := new(CloseResponse)
	if err := json.Unmarshal(body, closeResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, closeResponse, nil
}
