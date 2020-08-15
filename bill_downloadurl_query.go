// Package alipay https://docs.open.alipay.com/api_15/alipay.data.dataservice.bill.downloadurl.query
package alipay

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
)

const (
	// BillTypeTrade 指商户基于支付宝交易收单的业务账
	BillTypeTrade = "trade"

	// BillTypeSigncustomer 指基于商户支付宝余额收入及支出等资金变动的帐务账单
	BillTypeSigncustomer = "signcustomer"
)

// CsvBusinessType ...
const (
	CsvBusinessTypeTrade  = "交易"
	CsvBusinessTypeRefund = "退款"
)

// BillTradeEntry ...
type BillTradeEntry struct {
	TradeNo             string // 支付宝交易号
	OutTradeNo          string // 商户订单号
	BusinessType        string // 业务类型
	Subject             string // 商品名称
	TimeStart           string // 创建时间
	TimeEnd             string // 完成时间
	ShopNo              string // 门店编号
	ShopName            string // 门店名称
	Operator            string // 操作员
	TerminalNo          string // 终端号
	BuyerEmail          string // 对方账户
	TotalAmount         string // 订单金额（元）
	ReceiptAmount       string // 商家实收（元）
	Coupon              string // 支付宝红包（元）
	Jf                  string // 集分宝（元）
	AlipayOff           string // 支付宝优惠（元）
	SellerOff           string // 商家优惠（元）
	CouponChargeOff     string // 券核销金额（元）
	CouponName          string // 券名称
	SellerCouponConsume string // 商家红包消费金额（元）
	HandlingCharge      string // 卡消费金额（元）
	OutRequestNo        string // 退款批次号/请求号
	Service             string // 服务费（元）
	Fr                  string // 分润（元）
	Body                string // 备注
}

// BillSigncustomerEntry ...
type BillSigncustomerEntry struct {
	FundFlowID     string // 账务流水号
	TransactionID  string // 业务流水号
	BusinessID     string // 商户订单号
	ProductName    string // 商品名称
	TimeStart      string // 发生时间
	OtherAccount   string // 对方账号
	IncomeAmount   string // 收入金额（+元）
	ExpensesAmount string // 支出金额（-元）
	Balance        string // 账户余额（元）
	TradingChannel string // 交易渠道
	BusinessType   string // 业务类型
	Remark         string // 备注
}

// BillDownloadURLQueryParam ...
type BillDownloadURLQueryParam struct {
	BillType string `json:"bill_type,omitempty"` // 账单类型，商户通过接口或商户经开放平台授权后其所属服务商通过接口可以获取以下账单类型：trade、signcustomer；trade指商户基于支付宝交易收单的业务账单；signcustomer是指基于商户支付宝余额收入及支出等资金变动的帐务账单；
	BillDate string `json:"bill_date,omitempty"` // 账单时间：日账单格式为yyyy-MM-dd，月账单格式为yyyy-MM。
}

// BillDownloadURLQueryResponse ...
type BillDownloadURLQueryResponse struct {
	ResponseError
	BillDownloadURL string `json:"bill_download_url"` // 账单下载地址链接，获取连接后30秒后未下载，链接地址失效。
}

// IsBillNotExist ...
func (resp *BillDownloadURLQueryResponse) IsBillNotExist() bool {
	return resp.SubCode == "isp.bill_not_exist"
}

// BillDownloadurlQuery ...
func (alipay *Alipay) BillDownloadurlQuery(param *BillDownloadURLQueryParam) (int, *BillDownloadURLQueryResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		MethodAlipayDataDataserviceBillDownloadurlQuery,
	)
	if err != nil {
		return 0, nil, err
	}
	billDownloadURLQueryResponse := new(BillDownloadURLQueryResponse)
	if err := json.Unmarshal(body, billDownloadURLQueryResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, billDownloadURLQueryResponse, nil
}

// DownloadBill ...
func (alipay *Alipay) DownloadBill(billURL string) ([]byte, error) {
	resp, err := alipay.HTTPClient().Get(billURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// BillTradeList ...
func (alipay *Alipay) BillTradeList(bill []byte) ([]*BillTradeEntry, error) {
	billTradeEntryList := make([]*BillTradeEntry, 0)

	byteReader := bytes.NewReader(bill)

	zipReader, err := zip.NewReader(byteReader, int64(byteReader.Len()))
	if err != nil {
		return nil, err
	}

	for _, file := range zipReader.File {
		fileNameBytesInUtf8, _ := GbkToUtf8([]byte(file.Name))
		fileNameInUtf8 := string(fileNameBytesInUtf8)
		if strings.HasSuffix(fileNameInUtf8, "业务明细.csv") {

			var buf bytes.Buffer

			fileReaderCloser, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer fileReaderCloser.Close()

			if _, err := io.Copy(&buf, fileReaderCloser); err != nil {
				return nil, err
			}

			contentBytesInUtd8, _ := GbkToUtf8(buf.Bytes())

			content := string(contentBytesInUtd8)
			lines := strings.Split(content, "\n")
			var (
				validContent string
			)
			for _, line := range lines {
				if strings.HasPrefix(line, "#") {
					continue
				}
				if len(validContent) > 0 {
					validContent += "\n"
				}
				validContent += line
			}

			csvReader := csv.NewReader(strings.NewReader(validContent))
			records, err := csvReader.ReadAll()

			for _, record := range records {
				for idx, field := range record {
					record[idx] = strings.TrimSuffix(field, "\t")
				}
			}

			if len(records) > 0 {
				for _, record := range records[1:] {
					if len(record) != 25 {
						continue
					}

					entry := new(BillTradeEntry)
					billTradeEntryList = append(billTradeEntryList, entry)

					entry.TradeNo = record[0]
					entry.OutTradeNo = record[1]
					entry.BusinessType = record[2]
					entry.Subject = record[3]
					entry.TimeStart = record[4]
					entry.TimeEnd = record[5]
					entry.ShopNo = record[6]
					entry.ShopName = record[7]
					entry.Operator = record[8]
					entry.TerminalNo = record[9]
					entry.BuyerEmail = record[10]
					entry.TotalAmount = record[11]
					entry.ReceiptAmount = record[12]
					entry.Coupon = record[13]
					entry.Jf = record[14]
					entry.AlipayOff = record[15]
					entry.SellerOff = record[16]
					entry.CouponChargeOff = record[17]
					entry.CouponName = record[18]
					entry.SellerCouponConsume = record[19]
					entry.HandlingCharge = record[20]
					entry.OutRequestNo = record[21]
					entry.Service = record[22]
					entry.Fr = record[23]
					entry.Body = record[24]
				}
			}
			break
		}
	}
	return billTradeEntryList, nil
}

// BillSigncustomerList ...
func (alipay *Alipay) BillSigncustomerList(bill []byte) ([]*BillSigncustomerEntry, error) {
	billSigncustomerEntryList := make([]*BillSigncustomerEntry, 0)

	byteReader := bytes.NewReader(bill)

	zipReader, err := zip.NewReader(byteReader, int64(byteReader.Len()))
	if err != nil {
		return nil, err
	}

	for _, file := range zipReader.File {
		fileNameBytesInUtf8, _ := GbkToUtf8([]byte(file.Name))
		fileNameInUtf8 := string(fileNameBytesInUtf8)
		if strings.HasSuffix(fileNameInUtf8, "账务明细.csv") {

			var buf bytes.Buffer

			fileReaderCloser, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer fileReaderCloser.Close()

			if _, err := io.Copy(&buf, fileReaderCloser); err != nil {
				return nil, err
			}

			contentBytesInUtd8, _ := GbkToUtf8(buf.Bytes())

			content := string(contentBytesInUtd8)
			lines := strings.Split(content, "\n")
			var (
				validContent string
			)
			for _, line := range lines {
				if strings.HasPrefix(line, "#") {
					continue
				}
				if len(validContent) > 0 {
					validContent += "\n"
				}
				validContent += line
			}

			csvReader := csv.NewReader(strings.NewReader(validContent))
			records, err := csvReader.ReadAll()

			for _, record := range records {
				for idx, field := range record {
					record[idx] = strings.TrimSuffix(field, "\t")
				}
			}

			if len(records) > 0 {
				for _, record := range records[1:] {
					if len(record) != 12 {
						continue
					}

					entry := new(BillSigncustomerEntry)
					billSigncustomerEntryList = append(billSigncustomerEntryList, entry)

					entry.FundFlowID = record[0]
					entry.TransactionID = record[1]
					entry.BusinessID = record[2]
					entry.ProductName = record[3]
					entry.TimeStart = record[4]
					entry.OtherAccount = record[5]
					entry.IncomeAmount = record[6]
					entry.ExpensesAmount = record[7]
					entry.Balance = record[8]
					entry.TradingChannel = record[9]
					entry.BusinessType = record[10]
					entry.Remark = record[11]
				}
			}
			break
		}
	}
	return billSigncustomerEntryList, nil
}
