package main

import (
	"fmt"
	"log"

	"github.com/djpken/go-exc/types"
)

func main() {
	fmt.Println("=== Decimal.IsZero() Method Examples ===")

	// Example 1: Check if balance is zero
	balance := types.Decimal("0")
	if balance.IsZero() {
		fmt.Println("✓ Balance is zero, no funds available")
	}

	// Example 2: Check if price changed
	oldPrice := types.Decimal("100.50")
	newPrice := types.Decimal("100.50")
	priceChange := types.Decimal("0.0")

	if priceChange.IsZero() {
		fmt.Println("✓ Price unchanged:", oldPrice.String(), "→", newPrice.String())
	}

	// Example 3: Filter non-zero balances
	balances := map[string]types.Decimal{
		"BTC":  "0.5",
		"ETH":  "0",
		"USDT": "1000.00",
		"USDC": "0.00",
		"BNB":  "2.5",
	}

	fmt.Println("\n✓ Non-zero balances:")
	for currency, amount := range balances {
		if !amount.IsZero() {
			fmt.Printf("  %s: %s\n", currency, amount.String())
		}
	}

	// Example 4: Check minimum order quantity
	orderQty := types.Decimal("0.001")
	minQty := types.Decimal("0.01")

	if orderQty.IsZero() {
		log.Println("✗ Order quantity cannot be zero")
	} else {
		orderQtyFloat, _ := orderQty.Float64()
		minQtyFloat, _ := minQty.Float64()
		if orderQtyFloat < minQtyFloat {
			fmt.Printf("\n✗ Order quantity %s is below minimum %s\n", orderQty, minQty)
		}
	}

	// Example 5: Using ZeroDecimal constant
	const defaultFee = types.ZeroDecimal
	if defaultFee.IsZero() {
		fmt.Println("\n✓ Using zero fee (maker rebate)")
	}

	// Example 6: Various zero representations
	fmt.Println("\n=== Zero Representations ===")
	zeroValues := []types.Decimal{
		"0",
		"0.0",
		"0.00",
		"0.000000",
		"",
		types.ZeroDecimal,
	}

	for _, val := range zeroValues {
		fmt.Printf("Decimal(%q).IsZero() = %v\n", val, val.IsZero())
	}

	// Example 7: Non-zero values
	fmt.Println("\n=== Non-Zero Values ===")
	nonZeroValues := []types.Decimal{
		"1",
		"0.1",
		"-0.1",
		"1000.50",
	}

	for _, val := range nonZeroValues {
		fmt.Printf("Decimal(%q).IsZero() = %v\n", val, val.IsZero())
	}
}
