// Package alipay https://docs.open.alipay.com/api_1/alipay.trade.pay/
package alipay

import (
	"encoding/json"
)

// PayParamGoods ...
type PayParamGoods struct {
	GoodsID       string `json:"goods_id"`       // 商品的编号
	GoodsName     string `json:"goods_name"`     // 商品名称
	Quantity      int    `json:"quantity"`       // 商品数量
	Price         string `json:"price"`          // 商品单价
	GoodsCategory string `json:"goods_category"` // 商品类目
	Body          string `json:"body"`           // 商品描述信息
	ShowURL       string `json:"show_url"`       // 商品的展示地址
}

// PayParamExtendParams ...
type PayParamExtendParams struct {
	SysServiceProviderID string `json:"sys_service_provider_id"` // 系统商编号
}

// PayParam ...
type PayParam struct {
	OutTradeNo         string               `json:"out_trade_no"`        // 商户订单号
	Scene              string               `json:"scene"`               // 支付场景
	AuthCode           string               `json:"auth_code"`           // 支付授权码
	ProductCode        string               `json:"product_code"`        // 销售产品码
	Subject            string               `json:"subject"`             // 订单标题
	BuyerID            string               `json:"buyer_id"`            // 买家的支付宝用户id
	SellerID           string               `json:"seller_id "`          // 商户签约账号对应的支付宝用户ID
	TotalAmount        string               `json:"total_amount"`        // 订单总金额
	DiscountableAmount string               `json:"discountable_amount"` // 参与优惠计算的金额
	Body               string               `json:"body"`                // 订单描述
	GoodsDetailList    []*PayParamGoods     `json:" goods_detail"`       // 订单包含的商品列表信息
	OperatorID         string               `json:"operator_id"`         // 商户操作员编号
	StoreID            string               `json:"store_id"`            // 商户门店编号
	TerminalID         string               `json:"terminal_id"`         // 商户机具终端编号
	ExtendParams       PayParamExtendParams `json:"extend_params"`       // 业务扩展参数
	TimeoutExpress     string               `json:"timeout_express"`     // 该笔订单允许的最晚付款时间
}

// PayResponseFundBill ...
type PayResponseFundBill struct {
	FundChannel string `json:"fund_channel"` // 交易使用的资金渠道
	Amount      string `json:"amount"`       // 该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount"`  // 渠道实际付款金额
}

// PayResponseVoucherDetail ...
type PayResponseVoucherDetail struct {
	ID                         string `json:"id"`                           // 券id
	Name                       string `json:"name"`                         // 券名称
	Type                       string `json:"type"`                         // 当前有三种类型
	Amount                     string `json:"amount"`                       // 优惠券面额
	MerchantContribute         string `json:"merchant_contribute"`          // 商家出资
	OtherContribute            string `json:" other_contribute"`            // 其他出资方出资金额
	Meomo                      string `json:" memo"`                        // 优惠券备注信息
	TemplateID                 string `json:" template_id"`                 // 券模板id
	PurchaseBuyerContribute    string `json:"purchase_buyer_contribute"`    // 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时用户实际付款的金额
	PurchaseMerchantContribute string `json:"purchase_merchant_contribute"` // 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时商户优惠的金额
	PurchaseAntContribute      string `json:"purchase_ant_contribute"`      // 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时平台优惠的金额
}

// PayResponse ...
type PayResponse struct {
	ResponseError
	TradeNo             string                      `json:"trade_no"`              // 支付宝交易号
	OutTradeNo          string                      `json:"out_trade_no"`          // 商户订单号
	BuyerLogonID        string                      `json:"buyer_logon_id"`        // 买家支付宝账号
	TotalAmount         string                      `json:"total_amount"`          // 交易金额
	ReceiptAmount       string                      `json:"receipt_amount"`        // 实收金额
	BuyerPayAmount      string                      `json:"buyer_pay_amount"`      // 买家付款的金额
	PointAmount         string                      `json:"point_amount"`          // 使用积分宝付款的金额
	InvoiceAmount       string                      `json:"invoice_amount"`        // 交易中可给用户开具发票的金额
	GmtPayment          string                      `json:"gmt_payment"`           // 交易支付时间
	FundBillList        []*PayResponseFundBill      `json:"fund_bill_list"`        // 交易支付使用的资金渠道
	CardBalance         string                      `json:"card_balance"`          // 支付宝卡余额
	StoreName           string                      `json:"store_name"`            // 发生支付交易的商户门店名称
	BuyerUserID         string                      `json:"buyer_user_id"`         // 买家在支付宝的用户id
	DiscountGoodsDetail string                      `json:"discount_goods_detail"` // 本次交易支付所使用的单品券优惠的商品优惠信息
	VoucherDetailList   []*PayResponseVoucherDetail `json:" voucher_detail_list"`  // 本交易支付时使用的所有优惠券信息
	BusinessParams      string                      `json:"business_params"`       // 商户传入业务信息
	BuyerUserType       string                      `json:"buyer_user_type"`       // 买家用户类型
}

// Pay ...
func (alipay *Alipay) Pay(param *PayParam, notifyURL, appAuthToken string) (int, *PayResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		MethodAlipayTradePay,
		WithNotifyURL(notifyURL),
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
