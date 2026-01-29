package main

import (
	"context"
	"fmt"
	"log"

	"github.com/djpken/go-exc"
)

// ========== Configuration Constants ==========
// 将 API 配置和交易所类型定义为常数，方便管理

const (
	// Exchange Type - 交易所类型
	EXCHANGE_TYPE = exc.BitMart // 可以改为 exc.OKX, exc.OKXTest 等

	// API Keys - API 密钥配置
	API_KEY    = "f42ce865e6123a77b507659452384a2b48165991"
	SECRET_KEY = "72cca2bc92ee262f7e5398e3344b15638d43224a6fc5404567fc884f655a4b73"
	MEMO       = ".Vm3djpcl3gj94" // BitMart 需要，OKEx 不需要
	PASSPHRASE = ""               // OKEx 需要，BitMart 不需要

	// Trading Pair - 交易对
	SYMBOL = "BTC_USDT" // BitMart 使用 BTC_USDT，OKEx 使用 BTC-USDT
)

// 这个示例展示如何使用 Exchange 接口处理单一交易所
func main() {
	ctx := context.Background()

	log.Println("=================================================================")
	log.Println("              单一交易所示例")
	log.Println("=================================================================")
	log.Println()

	// 创建交易所客户端
	client, err := createExchangeClient(ctx)
	if err != nil {
		log.Fatalf("❌ 创建客户端失败: %v", err)
	}
	defer client.Close()

	log.Printf("✅ 成功连接到 %s 交易所", client.Name())
	log.Println()

	balance, err := client.GetBalance(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取账户余额失败: %v", err)
	}
	fmt.Sprint(balance)
	// 获取市场数据
	//getMarketData(ctx, client)

	// 查询账户信息
	// getAccountInfo(ctx, client)

	// 下单示例（需要有效的 API 密钥）
	// placeOrderExample(ctx, client)  // See bitmart_trading_example for complete trading examples

	log.Println()
	log.Println("=================================================================")
	log.Println("✅ 示例执行完成")
	log.Println("=================================================================")
}

// createExchangeClient 创建交易所客户端
func createExchangeClient(ctx context.Context) (exc.Exchange, error) {
	// 根据交易所类型构建配置
	config := exc.Config{
		APIKey:    API_KEY,
		SecretKey: SECRET_KEY,
	}

	// BitMart 需要 memo
	if EXCHANGE_TYPE == exc.BitMart || EXCHANGE_TYPE == exc.BitMartTest {
		config.Extra = map[string]interface{}{
			"memo": MEMO,
		}
	}

	// OKEx 需要 passphrase
	if EXCHANGE_TYPE == exc.OKX || EXCHANGE_TYPE == exc.OKXTest {
		config.Passphrase = PASSPHRASE
	}

	return exc.NewExchange(ctx, EXCHANGE_TYPE, config)
}

// getMarketData 获取市场数据
func getMarketData(ctx context.Context, client exc.Exchange) {
	log.Println("【获取市场数据】")
	log.Println("-----------------------------------------------------------------")

	// 1. 获取行情
	ticker, err := client.GetTicker(ctx, SYMBOL)
	if err != nil {
		log.Printf("  ❌ 获取行情失败: %v", err)
	} else {
		log.Printf("  📊 %s 最新价格: %s", SYMBOL, ticker.LastPrice)
		log.Printf("     24h 最高: %s", ticker.High24h)
		log.Printf("     24h 最低: %s", ticker.Low24h)
		log.Printf("     24h 成交量: %s", ticker.Volume24h)
	}

	log.Println()
}
