package alipay

import (
	"encoding/json"
)

// https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query/

type FastpayRefundQueryParam struct {
	// 支付宝交易号
	TradeNo string `json:"trade_no,omitempty"`
	// 户订单号
	OutTradeNo string `json:"out_trade_no,omitempty"`
	// 请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的外部交易号
	OutRequestNo string `json:"out_request_no,omitempty"`
}

type FastpayRefundQueryResponse struct {
	ResponseError
	// 支付宝交易号
	TradeNo string `json:"trade_no"`
	// 商户订单号
	OutTradeNo string `json:"out_trade_no"`
	// 本笔退款对应的退款请求号
	OutRequestNo string `json:"out_request_no"`
	// 发起退款时，传入的退款原因
	RefundReason string `json:"refund_reason"`
	// 该笔退款所对应的交易的订单金额
	TotalAmount string `json:"total_amount"`
	// 本次退款请求，对应的退款金额
	RefundAmount string `json:"refund_amount"`
}

// 商户可使用该接口查询自已通过alipay.trade.refund提交的退款请求是否执行成功。 该接口的返回码10000，仅代表本次查询操作成功，不代表退款成功。如果该接口返回了查询数据，则代表退款成功，如果没有查询到则代表未退款成功，可以调用退款接口进行重试。重试时请务必保证退款请求号一致。
func (resp *FastpayRefundQueryResponse) IsRefundSuccess() bool {
	return resp.Success() && Float64ifyPrice(resp.RefundAmount) > 0
}

func (resp *FastpayRefundQueryResponse) IsNeedRetry() bool {
	return resp.Success() && Float64ifyPrice(resp.RefundAmount) == 0
}

func (resp *FastpayRefundQueryResponse) IsTradeNotExist() bool {
	return resp.SubCode == "ACQ.TRADE_NOT_EXIST"
}

func (alipay *Alipay) FastpayRefundQuery(param *FastpayRefundQueryParam) (int, *FastpayRefundQueryResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_TRADE_FASTPAY_REFUND_QUERY,
	)
	if err != nil {
		return 0, nil, err
	}
	refundQueryResponse := new(FastpayRefundQueryResponse)
	if err := json.Unmarshal(body, refundQueryResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, refundQueryResponse, nil
}
