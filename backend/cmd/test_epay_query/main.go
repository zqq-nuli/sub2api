package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// EpayOrderQueryResponse 易支付订单查询响应
type EpayOrderQueryResponse struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	TradeNo    string `json:"trade_no"`
	OutTradeNo string `json:"out_trade_no"`
	Type       string `json:"type"`
	PID        string `json:"pid"`
	AddTime    string `json:"addtime"`
	EndTime    string `json:"endtime"`
	Name       string `json:"name"`
	Money      string `json:"money"`
	Status     int    `json:"status"`
}

func main() {
	// 从环境变量或命令行参数获取配置
	apiURL := os.Getenv("EPAY_API_URL")
	if apiURL == "" {
		apiURL = "https://credit.linux.do/epay"
	}

	merchantID := os.Getenv("EPAY_MERCHANT_ID")
	merchantKey := os.Getenv("EPAY_MERCHANT_KEY")

	if merchantID == "" || merchantKey == "" {
		fmt.Println("请设置环境变量 EPAY_MERCHANT_ID 和 EPAY_MERCHANT_KEY")
		fmt.Println("例如: EPAY_MERCHANT_ID=xxx EPAY_MERCHANT_KEY=yyy go run main.go <trade_no>")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("用法: go run main.go <trade_no> [out_trade_no]")
		fmt.Println("trade_no: 支付平台订单号（如 011449220309450752）")
		fmt.Println("out_trade_no: 商户订单号（可选）")
		os.Exit(1)
	}

	tradeNo := os.Args[1]
	outTradeNo := ""
	if len(os.Args) > 2 {
		outTradeNo = os.Args[2]
	}

	// 构建查询URL
	u, err := url.Parse(apiURL + "/api.php")
	if err != nil {
		fmt.Printf("URL解析失败: %v\n", err)
		os.Exit(1)
	}

	q := u.Query()
	q.Set("act", "order")
	q.Set("pid", merchantID)
	q.Set("key", merchantKey)
	q.Set("trade_no", tradeNo)
	if outTradeNo != "" {
		q.Set("out_trade_no", outTradeNo)
	}
	u.RawQuery = q.Encode()

	fmt.Printf("查询URL: %s\n", u.String())
	fmt.Println("正在查询...")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(u.String())
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n原始响应: %s\n\n", string(body))

	// 解析响应
	var result EpayOrderQueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("解析响应失败: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	fmt.Println("=== 查询结果 ===")
	fmt.Printf("状态码: %d\n", result.Code)
	fmt.Printf("消息: %s\n", result.Msg)
	fmt.Printf("支付平台订单号: %s\n", result.TradeNo)
	fmt.Printf("商户订单号: %s\n", result.OutTradeNo)
	fmt.Printf("支付类型: %s\n", result.Type)
	fmt.Printf("商户ID: %s\n", result.PID)
	fmt.Printf("创建时间: %s\n", result.AddTime)
	fmt.Printf("完成时间: %s\n", result.EndTime)
	fmt.Printf("订单名称: %s\n", result.Name)
	fmt.Printf("金额: %s\n", result.Money)
	fmt.Printf("支付状态: %d (1=成功, 0=处理中/失败)\n", result.Status)

	if result.Code == 1 && result.Status == 1 {
		fmt.Println("\n✅ 订单已支付成功！")
	} else if result.Code == 1 && result.Status == 0 {
		fmt.Println("\n⏳ 订单处理中或失败")
	} else {
		fmt.Println("\n❌ 查询失败或订单不存在")
	}
}
