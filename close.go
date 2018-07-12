package alipay

import (
	"encoding/json"
)

type CloseParam struct {
	// 支付宝交易号
	TradeNo string `json:"trade_no,omitempty"`
	// 商户订单号
	OutTradeNo string `json:"out_trade_no,omitempty"`
	// 卖家端自定义的的操作员 ID
	OperatorId string `json:"operator_id,omitempty"`
}

type CloseResponse struct {
	ResponseError
	// 支付宝交易号
	TradeNo string `json:"trade_no"`
	// 商户订单号
	OutTradeNo string `json:"out_trade_no"`
}

func (alipay *Alipay) Close(param *CloseParam) (int, *CloseResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_TRADE_CLOSE,
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
