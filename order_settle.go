package alipay

import (
	"encoding/json"
)

// OrderSettleParamOpenAPIRoyaltyDetailInfoPojo ...
type OrderSettleParamOpenAPIRoyaltyDetailInfoPojo struct {
	TransOut         string  `json:" trans_out"`        // 分账支出方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	TransIn          string  `json:"trans_in"`          // 分账收入方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	Amount           float64 `json:"amount"`            // 分账的金额，单位为元
	AmountPercentage float64 `json:"amount_percentage"` // 分账信息中分账百分比。取值范围为大于0，少于或等于100的整数。
	Desc             string  `json:"desc"`              // 分账描述
}

// OrderSettleParam ...
type OrderSettleParam struct {
	OutRequestNo      string                                          `json:"out_request_no"`     // 结算请求流水号
	TradeNo           string                                          `json:"trade_no"`           // 支付宝订单号
	RoyaltyParameters []*OrderSettleParamOpenAPIRoyaltyDetailInfoPojo `json:"royalty_parameters"` // 分账明细信息
	OperatorID        string                                          `json:"operator_id"`        // 操作员id
}

// OrderSettleResponse ...
type OrderSettleResponse struct {
	ResponseError
	TradeNo string `json:"trade_no"` // 支付宝交易号
}

// OrderSettle ...
func (alipay *Alipay) OrderSettle(param *OrderSettleParam) (int, *FastpayRefundQueryResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		MethodAlipayTradeOrderSettle,
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
