package alipay

import (
	"errors"
	"net/url"
)

// ErrAsyncVerify ...
var (
	ErrAsyncVerify = errors.New("支付宝异步验签失败")
)

// ParseAsyncResponse ...
func (alipay *Alipay) ParseAsyncResponse(values url.Values) (*AsyncResponse, error) {
	if err := alipay.asyncVerifyRequest(values); err != nil {
		return nil, ErrAsyncVerify
	}

	asyncResponse := new(AsyncResponse)

	asyncResponse.NotifyTime = values.Get("notify_time")
	asyncResponse.NotifyType = values.Get("notify_type")
	asyncResponse.NotifyID = values.Get("notify_id")
	asyncResponse.AppID = values.Get("app_id")
	asyncResponse.Charset = values.Get("charset")
	asyncResponse.Version = values.Get("version")
	asyncResponse.SignType = values.Get("sign_type")
	asyncResponse.Sign = values.Get("sign")
	asyncResponse.AuthAppID = values.Get("auth_app_id")
	asyncResponse.TradeNo = values.Get("trade_no")
	asyncResponse.OutTradeNo = values.Get("out_trade_no")
	asyncResponse.OutBizNo = values.Get("out_biz_no")
	asyncResponse.BuyerID = values.Get("buyer_id")
	asyncResponse.BuyerLogonID = values.Get("buyer_logon_id")
	asyncResponse.SellerID = values.Get("seller_id")
	asyncResponse.SellerEmail = values.Get("seller_email")
	asyncResponse.TradeStatus = values.Get("trade_status")
	asyncResponse.TotalAmount = values.Get("total_amount")
	asyncResponse.ReceiptAmount = values.Get("receipt_amount")
	asyncResponse.InvoiceAmount = values.Get("invoice_amount")
	asyncResponse.BuyerPayAmount = values.Get("buyer_pay_amount")
	asyncResponse.PointAmount = values.Get("point_amount")
	asyncResponse.RefundFee = values.Get("refund_fee")
	asyncResponse.Subject = values.Get("subject")
	asyncResponse.Body = values.Get("body")
	asyncResponse.GmtCreate = values.Get("gmt_create")
	asyncResponse.GmtPayment = values.Get("gmt_payment")
	asyncResponse.GmtRefund = values.Get("gmt_refund")
	asyncResponse.GmtClose = values.Get("gmt_close")
	asyncResponse.FundBillList = values.Get("fund_bill_list")
	asyncResponse.PassbackParams = values.Get("passback_params")
	asyncResponse.VoucherDetailList = values.Get("voucher_detail_list")

	return asyncResponse, nil
}

// AsyncResponse ...
type AsyncResponse struct {
	NotifyTime        string
	NotifyType        string
	NotifyID          string
	AppID             string
	Charset           string
	Version           string
	SignType          string
	Sign              string
	AuthAppID         string
	TradeNo           string
	OutTradeNo        string
	OutBizNo          string
	BuyerID           string
	BuyerLogonID      string
	SellerID          string
	SellerEmail       string
	TradeStatus       string
	TotalAmount       string
	ReceiptAmount     string
	InvoiceAmount     string
	BuyerPayAmount    string
	PointAmount       string
	RefundFee         string
	Subject           string
	Body              string
	GmtCreate         string
	GmtPayment        string
	GmtRefund         string
	GmtClose          string
	FundBillList      string
	PassbackParams    string
	VoucherDetailList string
}
