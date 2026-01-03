package service

import (
	"crypto/md5"
	cryptorand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// EpayClient 易支付客户端
type EpayClient struct {
	apiURL      string
	merchantID  string
	merchantKey string
}

// NewEpayClient 创建易支付客户端
func NewEpayClient(apiURL, merchantID, merchantKey string) *EpayClient {
	return &EpayClient{
		apiURL:      apiURL,
		merchantID:  merchantID,
		merchantKey: merchantKey,
	}
}

// CreatePaymentURL 构建易支付支付URL
func (c *EpayClient) CreatePaymentURL(params map[string]string) (string, error) {
	// 添加商户ID
	params["pid"] = c.merchantID

	// 生成签名
	sign := c.GenerateSign(params)
	params["sign"] = sign
	params["sign_type"] = "MD5"

	// 构建URL (linux.do credit API: /pay/submit.php)
	u, err := url.Parse(c.apiURL + "/pay/submit.php")
	if err != nil {
		return "", fmt.Errorf("invalid api url: %w", err)
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// CreatePaymentForm 构建易支付表单数据（用于POST提交）
// 返回：表单提交URL、表单参数（含签名）、错误
func (c *EpayClient) CreatePaymentForm(params map[string]string) (string, map[string]string, error) {
	// 添加商户ID
	params["pid"] = c.merchantID

	// 生成签名
	sign := c.GenerateSign(params)
	params["sign"] = sign
	params["sign_type"] = "MD5"

	// 构建提交URL (linux.do credit API: /pay/submit.php)
	submitURL := c.apiURL + "/pay/submit.php"

	return submitURL, params, nil
}

// VerifySign 验证易支付回调签名
func (c *EpayClient) VerifySign(params map[string]string) bool {
	receivedSign := params["sign"]
	if receivedSign == "" {
		return false
	}

	expectedSign := c.GenerateSign(params)
	return receivedSign == expectedSign
}

// GenerateSign 生成易支付签名
// 签名规则：
// 1. 筛选：获取所有请求参数，不包括sign、sign_type字段
// 2. 排序：按照参数名ASCII码从小到大排序
// 3. 拼接：按照 key=value&key=value 的格式拼接
// 4. 追加密钥：在拼接串最后追加商户密钥
// 5. MD5：对拼接串进行MD5加密，并转小写
func (c *EpayClient) GenerateSign(params map[string]string) string {
	// 1. 筛选参数
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" && k != "sign_type" && params[k] != "" {
			keys = append(keys, k)
		}
	}

	// 2. 排序
	sort.Strings(keys)

	// 3. 拼接
	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			_, _ = sb.WriteString("&")
		}
		_, _ = sb.WriteString(k)
		_, _ = sb.WriteString("=")
		_, _ = sb.WriteString(params[k])
	}

	// 4. 追加密钥
	_, _ = sb.WriteString(c.merchantKey)

	// 5. MD5加密
	hash := md5.Sum([]byte(sb.String()))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}

// GenerateOrderNo 生成订单号：ORD{timestamp}{random12}
// 使用 crypto/rand 生成密码学安全的随机数
// 总长度: 3 + 14 + 12 = 29 chars，符合 VARCHAR(32) 限制
func GenerateOrderNo() string {
	timestamp := time.Now().Format("20060102150405")
	random := generateSecureRandomString(12)
	return "ORD" + timestamp + random
}

// generateSecureRandomString 生成密码学安全的随机字符串
func generateSecureRandomString(length int) string {
	// 使用更长的随机字节数以增加安全性
	randomBytes := make([]byte, length)
	if _, err := cryptorand.Read(randomBytes); err != nil {
		// 极端情况下回退到时间戳+进程ID（虽然不安全但至少唯一）
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return strings.ToUpper(hex.EncodeToString(randomBytes)[:length])
}

// IsValidPaymentMethod 验证支付方式是否有效
func IsValidPaymentMethod(method string) bool {
	validMethods := map[string]bool{
		"alipay": true,
		"wxpay":  true,
		"qqpay":  true,
		"epusdt": true,
	}
	return validMethods[method]
}

// EpayOrderQueryResponse 易支付订单查询响应
type EpayOrderQueryResponse struct {
	Code       int    `json:"code"`        // 1=成功，-1=失败
	Msg        string `json:"msg"`         // 消息
	TradeNo    string `json:"trade_no"`    // 支付平台订单号
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	Type       string `json:"type"`        // 支付类型
	PID        string `json:"pid"`         // 商户ID
	AddTime    string `json:"addtime"`     // 创建时间
	EndTime    string `json:"endtime"`     // 完成时间
	Name       string `json:"name"`        // 订单名称
	Money      string `json:"money"`       // 金额
	Status     int    `json:"status"`      // 状态：1=成功，0=处理中/失败
}

// QueryOrderStatus 查询上游支付订单状态
// trade_no: 支付平台订单号（必填）
// out_trade_no: 商户订单号（可选）
func (c *EpayClient) QueryOrderStatus(tradeNo, outTradeNo string) (*EpayOrderQueryResponse, error) {
	// 构建查询URL
	u, err := url.Parse(c.apiURL + "/api.php")
	if err != nil {
		return nil, fmt.Errorf("invalid api url: %w", err)
	}

	q := u.Query()
	q.Set("act", "order")
	q.Set("pid", c.merchantID)
	q.Set("key", c.merchantKey)
	q.Set("trade_no", tradeNo)
	if outTradeNo != "" {
		q.Set("out_trade_no", outTradeNo)
	}
	u.RawQuery = q.Encode()

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("query order failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// 解析响应
	var result EpayOrderQueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response failed: %w, body: %s", err, string(body))
	}

	return &result, nil
}
