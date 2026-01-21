package main

import (
	"context"
	"errors"
	"log"

	"github.com/djpken/go-exc"
)

// 这个示例展示如何使用 GetConfig 方法
// 以及如何处理不支持的功能

func main() {
	ctx := context.Background()

	log.Println("=================================================================")
	log.Println("              GetConfig 方法示例")
	log.Println("=================================================================")
	log.Println()

	// ========== BitMart 示例（不支持）==========
	log.Println("【BitMart 交易所】")
	log.Println("-----------------------------------------------------------------")
	testBitMartConfig(ctx)

	log.Println()

	// ========== OKEx 示例（支持）==========
	log.Println("【OKEx 交易所】")
	log.Println("-----------------------------------------------------------------")
	testOKExConfig(ctx)

	log.Println()
	log.Println("=================================================================")
	log.Println("✅ 示例执行完成！")
	log.Println("=================================================================")
}

// testBitMartConfig 测试 BitMart GetConfig（不支持）
func testBitMartConfig(ctx context.Context) {
	// 创建 BitMart 客户端
	client, err := exc.NewExchange(ctx, exc.BitMart, exc.Config{
		APIKey:    "f42ce865e6123a77b507659452384a2b48165991",
		SecretKey: "72cca2bc92ee262f7e5398e3344b15638d43224a6fc5404567fc884f655a4b73",
		Extra: map[string]interface{}{
			"memo": ".Vm3djpcl3gj94",
		},
	})
	if err != nil {
		log.Printf("❌ 创建客户端失败: %v", err)
		return
	}
	defer client.Close()

	log.Printf("✅ 创建成功: %s", client.Name())

	// 尝试获取配置
	config, err := client.GetConfig(ctx)

	// 检查错误类型
	if errors.Is(err, exc.ErrNotSupported) {
		log.Printf("  ⚠️  GetConfig 不支持: %v", err)
		log.Printf("  ℹ️  BitMart 交易所不提供账户配置 API")
	} else if err != nil {
		log.Printf("  ❌ 获取配置失败: %v", err)
	} else if config != nil {
		log.Printf("  ✅ 获取配置成功: %+v", config)
	}
}

// testOKExConfig 测试 OKEx GetConfig（支持）
func testOKExConfig(ctx context.Context) {
	// 创建 OKEx 客户端
	client, err := exc.NewExchange(ctx, exc.OKXTest, exc.Config{
		APIKey:     "your-api-key",
		SecretKey:  "your-secret-key",
		Passphrase: "your-passphrase",
	})
	if err != nil {
		log.Printf("❌ 创建客户端失败: %v", err)
		return
	}
	defer client.Close()

	log.Printf("✅ 创建成功: %s", client.Name())

	// 获取配置
	config, err := client.GetConfig(ctx)
	if err != nil {
		if errors.Is(err, exc.ErrNotSupported) {
			log.Printf("  ⚠️  GetConfig 不支持: %v", err)
		} else {
			log.Printf("  ⚠️  获取配置失败（可能是测试密钥）: %v", err)
		}
		return
	}

	if config != nil {
		log.Printf("  ✅ 账户配置:")
		log.Printf("     UID: %s", config.UID)
		log.Printf("     Level: %s", config.Level)
		log.Printf("     AutoLoan: %v", config.AutoLoan)
		log.Printf("     PositionMode: %s", config.PositionMode)
		if len(config.Extra) > 0 {
			log.Printf("     Extra: %+v", config.Extra)
		}
	}
}

// demonstrateGenericHandler 展示如何编写通用的处理函数
func demonstrateGenericHandler(ctx context.Context, exchange exc.Exchange) {
	log.Printf("处理 %s...", exchange.Name())

	// 尝试获取配置
	config, err := exchange.GetConfig(ctx)

	// 优雅地处理不支持的情况
	if errors.Is(err, exc.ErrNotSupported) {
		log.Printf("  ⚠️  %s 不支持 GetConfig", exchange.Name())
		// 继续执行其他操作...
		return
	}

	if err != nil {
		log.Printf("  ❌ 获取配置失败: %v", err)
		return
	}

	if config != nil {
		log.Printf("  ✅ UID: %s, Level: %s", config.UID, config.Level)
	}
}

// 展示批量处理多个交易所
func batchProcessExchanges(ctx context.Context) {
	exchanges := []exc.Exchange{}

	// 添加多个交易所
	bitmart, _ := exc.NewExchange(ctx, exc.BitMart, exc.Config{
		APIKey:    "key",
		SecretKey: "secret",
		Extra:     map[string]interface{}{"memo": "memo"},
	})
	if bitmart != nil {
		exchanges = append(exchanges, bitmart)
		defer bitmart.Close()
	}

	okex, _ := exc.NewExchange(ctx, exc.OKXTest, exc.Config{
		APIKey:     "key",
		SecretKey:  "secret",
		Passphrase: "pass",
	})
	if okex != nil {
		exchanges = append(exchanges, okex)
		defer okex.Close()
	}

	// 对所有交易所执行相同操作
	for _, ex := range exchanges {
		demonstrateGenericHandler(ctx, ex)
	}
}
