package alipay

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// https://docs.open.alipay.com/api_15/alipay.data.dataservice.bill.downloadurl.query

const (
	// 指商户基于支付宝交易收单的业务账
	BILL_TYPE_TRADE = "trade"

	// 指基于商户支付宝余额收入及支出等资金变动的帐务账单
	BILL_TYPE_SIGNCUSTOMER = "signcustomer"
)

const (
	CSV_BUSINESS_TYPE_TRADE  = "交易"
	CSV_BUSINESS_TYPE_REFUND = "退款"
)

type BillTradeEntry struct {
	// 支付宝交易号
	TradeNo string
	// 商户订单号
	OutTradeNo string
	// 业务类型
	BusinessType string
	// 商品名称
	Subject string
	// 创建时间
	TimeStart string
	// 完成时间
	TimeEnd string
	// 门店编号
	ShopNo string
	// 门店名称
	ShopName string
	// 操作员
	Operator string
	// 终端号
	TerminalNo string
	// 对方账户
	BuyerEmail string
	// 订单金额（元）
	TotalAmount string
	// 商家实收（元）
	ReceiptAmount string
	// 支付宝红包（元）
	Coupon string
	// 集分宝（元）
	Jf string
	// 支付宝优惠（元）
	AlipayOff string
	// 商家优惠（元）
	SellerOff string
	// 券核销金额（元）
	CouponChargeOff string
	// 券名称
	CouponName string
	// 商家红包消费金额（元）
	SellerCouponConsume string
	// 卡消费金额（元）
	HandlingCharge string
	// 退款批次号/请求号
	OutRequestNo string
	// 服务费（元）
	Service string
	// 分润（元）
	Fr string
	// 备注
	Body string
}

type BillDownloadUrlQueryParam struct {
	// 账单类型，商户通过接口或商户经开放平台授权后其所属服务商通过接口可以获取以下账单类型：trade、signcustomer；trade指商户基于支付宝交易收单的业务账单；signcustomer是指基于商户支付宝余额收入及支出等资金变动的帐务账单；
	BillType string `json:"bill_type,omitempty"`
	// 账单时间：日账单格式为yyyy-MM-dd，月账单格式为yyyy-MM。
	BillDate string `json:"bill_date,omitempty"`
}

type BillDownloadUrlQueryResponse struct {
	ResponseError
	// 账单下载地址链接，获取连接后30秒后未下载，链接地址失效。
	BillDownloadUrl string `json:"bill_download_url"`
}

func (resp *BillDownloadUrlQueryResponse) IsBillNotExist() bool {
	return resp.SubCode == "isp.bill_not_exist"
}

func (alipay *Alipay) BillDownloadurlQuery(param *BillDownloadUrlQueryParam) (int, *BillDownloadUrlQueryResponse, error) {
	statusCode, body, err := alipay.OnRequest(
		param,
		METHOD_ALIPAY_DATA_DATASERVICE_BILL_DOWNLOADURL_QUERY,
	)
	if err != nil {
		return 0, nil, err
	}
	billDownloadUrlQueryResponse := new(BillDownloadUrlQueryResponse)
	if err := json.Unmarshal(body, billDownloadUrlQueryResponse); err != nil {
		return 0, nil, err
	}
	return statusCode, billDownloadUrlQueryResponse, nil
}

func (alipay *Alipay) DownloadBill(billUrl string) ([]byte, error) {
	resp, err := alipay.HttpClient().Get(billUrl)
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

			fileReadCloser, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer fileReadCloser.Close()

			var buf bytes.Buffer

			if _, err := io.Copy(&buf, fileReadCloser); err != nil {
				return nil, err
			}

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

			log.Println(content)

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
