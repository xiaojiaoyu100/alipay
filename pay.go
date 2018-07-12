package alipay

import (
	"encoding/json"
)

// https://docs.open.alipay.com/api_1/alipay.trade.pay/

type PayParamGoods struct {
	// 商品的编号
	GoodsId string `json:"goods_id"`
	// 商品名称
	GoodsName string `json:"goods_name"`
	// 商品数量
	Quantity int `json:"quantity"`
	// 商品单价
	Price string `json:"price"`
	// 商品类目
	GoodsCategory string `json:"goods_category"`
	// 商品描述信息
	Body string `json:"body"`
	// 商品的展示地址
	ShowUrl string `json:"show_url"`
}

type PayParamExtendParams struct {
	// 系统商编号
	SysServiceProviderId string `json:"sys_service_provider_id"`
}

type PayParam struct {
	// 商户订单号
	OutTradeNo string `json:"out_trade_no"`
	// 支付场景
	Scene string `json:"scene"`
	// 支付授权码
	AuthCode string `json:"auth_code"`
	// 销售产品码
	ProductCode string `json:"product_code "`
	// 订单标题
	Subject string `json:"subject"`
	// 买家的支付宝用户id
	BuyerId string `json:"buyer_id"`
	// 商户签约账号对应的支付宝用户ID
	SellerId string `json:"seller_id "`
	// 订单总金额
	TotalAmount string `json:"total_amount"`
	// 参与优惠计算的金额
	DiscountableAmount string `json:"discountable_amount"`
	// 订单描述
	Body string `json:"body"`
	// 订单包含的商品列表信息
	GoodsDetailList []*PayParamGoods `json:" goods_detail"`
	// 商户操作员编号
	OperatorId string `json:"operator_id"`
	// 商户门店编号
	StoreId string `json:"store_id"`
	// 商户机具终端编号
	TerminalId string `json:"terminal_id"`
	// 业务扩展参数
	ExtendParams PayParamExtendParams `json:"extend_params"`
	// 该笔订单允许的最晚付款时间
	TimeoutExpress string `json:"timeout_express"`
}

type PayResponseFundBill struct {
	// 交易使用的资金渠道
	FundChannel string `json:"fund_channel"`
	// 该支付工具类型所使用的金额
	Amount string `json:"amount"`
	// 渠道实际付款金额
	RealAmount string `json:"real_amount"`
}

type PayResponseVoucherDetail struct {
	// 券id
	Id string `json:"id"`
	// 券名称
	Name string `json:"name"`
	// 当前有三种类型
	Type string `json:"type"`
	// 优惠券面额
	Amount string `json:"amount"`
	// 商家出资
	MerchantContribute string `json:"merchant_contribute"`
	// 其他出资方出资金额
	OtherContribute string `json:" other_contribute"`
	// 优惠券备注信息
	Meomo string `json:" memo"`
	// 券模板id
	TemplateId string `json:" template_id"`
	// 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时用户实际付款的金额
	PurchaseBuyerContribute string `json:"purchase_buyer_contribute"`
	// 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时商户优惠的金额
	PurchaseMerchantContribute string `json:"purchase_merchant_contribute"`
	// 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时平台优惠的金额
	PurchaseAntContribute string `json:"purchase_ant_contribute"`
}

type PayResponse struct {
	ResponseError
	// 支付宝交易号
	TradeNo string `json:"trade_no"`
	// 商户订单号
	OutTradeNo string `json:"out_trade_no"`
	// 买家支付宝账号
	BuyerLogonId string `json:"buyer_logon_id"`
	// 交易金额
	TotalAmount string `json:"total_amount"`
	// 实收金额
	ReceiptAmount string `json:"receipt_amount"`
	// 买家付款的金额
	BuyerPayAmount string `json:"buyer_pay_amount"`
	// 使用积分宝付款的金额
	PointAmount string `json:"point_amount"`
	// 交易中可给用户开具发票的金额
	InvoiceAmount string `json:"invoice_amount"`
	// 交易支付时间
	GmtPayment string `json:"gmt_payment"`
	// 交易支付使用的资金渠道
	FundBillList []*PayResponseFundBill `json:"fund_bill_list"`
	// 支付宝卡余额
	CardBalance string `json:"card_balance"`
	// 发生支付交易的商户门店名称
	StoreName string `json:"store_name"`
	// 买家在支付宝的用户id
	BuyerUserId string `json:"buyer_user_id"`
	// 本次交易支付所使用的单品券优惠的商品优惠信息
	DiscountGoodsDetail string `json:"discount_goods_detail"`
	// 本交易支付时使用的所有优惠券信息
	VoucherDetailList []*PayResponseVoucherDetail `json:" voucher_detail_list"`
	// 商户传入业务信息
	BusinessParams string `json:"business_params"`
	// 买家用户类型
	BuyerUserType string `json:"buyer_user_type"`
}

func (alipay *Alipay) Pay(param *PayParam, notifyUrl, appAuthToken string) (int, *PayResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_TRADE_PAY,
		WithNotifyUrl(notifyUrl),
		WithAppAuthToken(appAuthToken),
	)
	if err != nil {
		return 0, nil, err
	}
	payResponse := new(PayResponse)
	if err := json.Unmarshal(body, payResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, payResponse, nil
}
