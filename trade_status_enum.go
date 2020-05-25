package alipay

const (
	// WaitBuyerPay 交易创建，等待买家付款
	WaitBuyerPay = "WAIT_BUYER_PAY"

	// TradeClosed 未付款交易超时关闭，或支付完成后全额退款
	TradeClosed = "TRADE_CLOSED"

	// TradeSuccess 交易支付成功
	TradeSuccess = "TRADE_SUCCESS"

	// TradeFinished 交易结束，不可退款
	TradeFinished = "TRADE_FINISHED"
)
