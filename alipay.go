// Package alipay provides alipay related api. See details on https://docs.open.alipay.com/api_1
package alipay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// AlipayGateway ...
const (
	AlipayGateway = "https://openapi.alipay.com/gateway.do"
)

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
	Timeout: time.Second * 60,
}

// CommonParam ...
type CommonParam struct {
	AppID        string `url:"app_id,omitempty"`         // 支付宝分配给开发者的应用ID
	Method       string `url:"method,omitempty"`         // 接口名称
	Format       string `url:"format,omitempty"`         // 仅支持JSON
	ReturnURL    string `url:"return_url,omitempty"`     // 同步返回地址
	Charset      string `url:"charset,omitempty"`        // 请求使用的编码格式，如utf-8,gbk,gb2312等
	SignType     string `url:"sign_type,omitempty"`      // 商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign         string `url:"sign,omitempty"`           // 商户请求参数的签名串
	Timestamp    string `url:"timestamp,omitempty"`      // 发送请求的时间，格式"yyyy-MM-dd HH:mm:ss"
	Version      string `url:"version,omitempty"`        // 调用的接口版本，固定为：1.0
	NotifyURL    string `url:"notify_url,omitempty"`     // 支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	AppAuthToken string `url:"app_auth_token,omitempty"` // 详见应用授权概述
	BizContent   string `url:"biz_content,omitempty"`    // 请求参数的集合
}

var defaultCommonParam = CommonParam{
	Format:   "JSON",
	Charset:  "utf-8",
	SignType: "RSA2",
	Version:  "1.0",
}

// Alipay ...
type Alipay struct {
	appID string
	*rsa.PrivateKey
	*rsa.PublicKey
	client *http.Client
}

// HTTPClient ...
func (alipay *Alipay) HTTPClient() *http.Client {
	return alipay.client
}

// AppID ...
func (alipay *Alipay) AppID() string {
	return alipay.appID
}

// Fill ...
type Fill func(*CommonParam)

// WithNotifyURL ...
func WithNotifyURL(notifyURL string) Fill {
	return func(cp *CommonParam) {
		cp.NotifyURL = notifyURL
	}
}

// WithReturnURL ...
func WithReturnURL(returnURL string) Fill {
	return func(cp *CommonParam) {
		cp.ReturnURL = returnURL
	}
}

// WithAppAuthToken ...
func WithAppAuthToken(authToken string) Fill {
	return func(cp *CommonParam) {
		cp.AppAuthToken = authToken
	}
}

// WithSignType ...
func WithSignType(signType string) Fill {
	return func(cp *CommonParam) {
		cp.SignType = signType
	}
}

// MakeParam ...
func (alipay *Alipay) MakeParam(content interface{}, method string, fillList ...Fill) (string, error) {
	biz, err := json.Marshal(&content)
	if err != nil {
		return "", fmt.Errorf("支付宝请求业务参数序列化失败: %w", err)
	}

	requestParam := CommonParam{
		AppID:      alipay.appID,
		Method:     method,
		Format:     defaultCommonParam.Format,
		Charset:    defaultCommonParam.Charset,
		SignType:   defaultCommonParam.SignType,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Version:    defaultCommonParam.Version,
		BizContent: string(biz),
	}

	for _, fill := range fillList {
		fill(&requestParam)
	}

	requestSignValues, err := query.Values(&requestParam)
	if err != nil {
		return "", fmt.Errorf("支付宝请求参数序列化失败: %w", err)
	}

	requestSignedStr := NormValues(requestSignValues)

	requestSign, err := RSA2(alipay.PrivateKey, requestSignedStr)
	if err != nil {
		return "", fmt.Errorf("支付宝商户私钥签名不成功: %w", err)
	}

	requestParam.Sign = requestSign
	values, err := query.Values(&requestParam)
	if err != nil {
		return "", fmt.Errorf("支付宝请求结构体不能序列化: %w", err)
	}

	return values.Encode(), nil
}

// OnRequest ...
func (alipay *Alipay) OnRequest(content interface{}, method string, fillList ...Fill) (int, []byte, error) {
	request, err := http.NewRequest(http.MethodGet, AlipayGateway, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("支付宝生成请求失败: %w", err)
	}

	param, err := alipay.MakeParam(content, method, fillList...)
	if err != nil {
		return 0, nil, fmt.Errorf("支付宝构造请求参数失败: %w", err)
	}

	request.URL.RawQuery = param
	response, err := alipay.client.Do(request)
	if err != nil {
		return 0, nil, fmt.Errorf("支付宝发起请求失败: %w", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("支付宝请求结果读取body失败: %w", err)
	}

	obj := make(map[string]*json.RawMessage)

	if err := json.Unmarshal(body, &obj); err != nil {
		return 0, nil, fmt.Errorf("支付宝请求结果反序列化失败: %w", err)
	}

	var (
		data []byte
		sig  []byte
	)

	for k, v := range obj {
		if strings.Contains(k, "response") {
			data = []byte(*v)
			break
		}
	}

	if rawMessage, ok := obj["sign"]; ok {
		sig = []byte(*rawMessage)
	}

	if len(sig) > 0 {
		var base64SignStr string

		if err := json.Unmarshal(sig, &base64SignStr); err != nil {
			return 0, nil, fmt.Errorf("反序列化签名失败: %w", err)
		}

		signStr, err := base64.StdEncoding.DecodeString(base64SignStr)
		if err != nil {
			return 0, nil, fmt.Errorf("base64解码签名失败: %w", err)
		}

		if err := Verify(alipay.PublicKey, data, signStr); err != nil {
			return 0, nil, fmt.Errorf("支付宝同步请求签名验证不通过: %w", err)
		}
	}

	return response.StatusCode, data, nil
}

func (alipay *Alipay) asyncVerifyRequest(values url.Values) error {
	// 验证参数步骤
	// 1) 去掉sign, sign_type
	// 2） 对剩下参数进行url decode
	// 3） 字典序排序
	// 4) sign解码
	base64Sign := values.Get("sign")

	toVerifyValues := url.Values{}

	for key, list := range values {
		if key == "sign" || key == "sign_type" {
			continue
		}

		for _, value := range list {
			toVerifyValues.Add(key, value)
		}
	}

	normToVerifyStr := NormValues(toVerifyValues)
	signBytes, _ := base64.StdEncoding.DecodeString(base64Sign)
	return Verify(alipay.PublicKey, []byte(normToVerifyStr), signBytes)
}

// NormValues ...
func NormValues(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k, l := range v {
		if len(l) == 0 {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			if v == "" {
				continue
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}

// New ...
func New(appID, publicKeyPath, privateKeyPath string) (*Alipay, error) {
	alipay := Alipay{
		appID:  appID,
		client: client,
	}

	var err error

	alipay.PublicKey, err = NewPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("支付宝公钥构建失败: %w", err)
	}

	alipay.PrivateKey, err = NewPrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("商户私钥构建失败: %w", err)
	}

	return &alipay, nil
}

// NewPublicKey publicKey参照ParsePKIXPublicKey
func NewPublicKey(path string) (pub *rsa.PublicKey, err error) {
	pubBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pubBlock, _ := pem.Decode(pubBytes)
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key format error")
	}
	return publicKey, nil
}

// NewPrivateKey privateKey格式参照ParsePKCS1PrivateKey
func NewPrivateKey(path string) (priKey *rsa.PrivateKey, err error) {
	privateBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	priBlock, _ := pem.Decode(privateBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(priBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// Verify ...
func Verify(publicKey *rsa.PublicKey, data []byte, sig []byte) error {
	h := sha256.New()
	h.Write(data)
	digest := h.Sum(nil)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, sig)
}

// RSA2 ...
func RSA2(privateKey *rsa.PrivateKey, in string) (string, error) {
	hash := sha256.New()
	_, err := io.WriteString(hash, in)
	if err != nil {
		return "", err
	}
	hashed := hash.Sum(nil)

	signBytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signBytes), nil
}
