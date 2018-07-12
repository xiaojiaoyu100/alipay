package alipay

import (
	"encoding/json"
)

// https://docs.open.alipay.com/api_1/alipay.trade.cancel/

type CancelParam struct {
	// 商户订单号
	TradeNo string `json:"trade_no"`
	// 支付宝交易号
	OutTradeNo string `json:"out_trade_no"`
}

type CancelResponse struct {
	ResponseError
	// 支付宝交易号
	TradeNo string `json:"trade_no"`
	// 商户订单号
	OutTradeNo string `json:"out_trade_no"`
	// 是否需要重试
	RetryFlag string `json:"retry_flag"`
	// 本次撤销触发的交易动作
	Action string `json:"action"`
}

func (alipay *Alipay) Cancel(param CancelParam) (int, *CancelResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_TRADE_CANCEL,
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
