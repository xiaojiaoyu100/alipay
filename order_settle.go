package alipay

import (
	"encoding/json"
)

type OrderSettleParamOpenApiRoyaltyDetailInfoPojo struct {
	// 分账支出方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	TransOut string `json:" trans_out"`

	// 分账收入方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	TransIn string `json:"trans_in"`

	// 分账的金额，单位为元
	Amount float64 `json:"amount"`

	// 分账信息中分账百分比。取值范围为大于0，少于或等于100的整数。
	AmountPercentage float64 `json:"amount_percentage"`

	// 分账描述
	Desc string `json:"desc"`
}

type OrderSettleParam struct {
	// 结算请求流水号
	OutRequestNo string `json:"out_request_no"`
	// 支付宝订单号
	TradeNo string `json:"trade_no"`
	// 分账明细信息
	RoyaltyParameters []*OrderSettleParamOpenApiRoyaltyDetailInfoPojo `json:"royalty_parameters"`
	// 操作员id
	OperatorId string `json:"operator_id"`
}

type OrderSettleResponse struct {
	ResponseError
	// 支付宝交易号
	TradeNo string `json:"trade_no"`
}

func (alipay *Alipay) OrderSettle(param *OrderSettleParam) (int, *FastpayRefundQueryResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_TRADE_ORDER_SETTLE,
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
