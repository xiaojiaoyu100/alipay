package alipay

// ResponseError ...
type ResponseError struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

// Success 这个判断不代表业务成功！
func (err ResponseError) Success() bool {
	return err.Code == GatewaySuccess && err.SubCode == ""
}
