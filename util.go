package alipay

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// StringifyPrice ...
func StringifyPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}

// Float64ifyPrice ...
func Float64ifyPrice(price string) float64 {
	p, _ := strconv.ParseFloat(price, 64)
	return p
}

// DayFormat ...
func DayFormat(t time.Time) string {
	return t.Format("20060102")
}

// GbkToUtf8 ...
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	changed, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return changed, nil
}
