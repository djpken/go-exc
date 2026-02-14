package main

import (
	"fmt"
	"log"

	"github.com/djpken/go-exc/types"
)

func main() {
	fmt.Println("=== Decimal IsZero() Examples ===\n")

	// Example 1: Check if price is zero before division
	price := types.MustDecimal("50000.25")
	quantity := types.MustDecimal("2")

	if price.IsZero() {
		log.Println("✗ Price cannot be zero")
	} else {
		total, _ := price.Mul(quantity)
		fmt.Printf("✓ Order total: %s × %s = %s USDT\n", price, quantity, total)
	}

	// Example 2: Detect price changes
	oldPrice := types.MustDecimal("100.50")
	newPrice := types.MustDecimal("100.50")
	priceChange := types.MustDecimal("0.0")

	if priceChange.IsZero() {
		fmt.Println("✓ Price unchanged:", oldPrice.String(), "→", newPrice.String())
	}

	// Example 3: Filter non-zero balances
	balances := map[string]types.Decimal{
		"BTC":  types.MustDecimal("0.5"),
		"ETH":  types.ZeroDecimal,
		"USDT": types.MustDecimal("1000.00"),
		"USDC": types.MustDecimal("0.00"),
		"BNB":  types.MustDecimal("2.5"),
	}

	fmt.Println("\n✓ Non-zero balances:")
	for currency, amount := range balances {
		if !amount.IsZero() {
			fmt.Printf("  %s: %s\n", currency, amount.String())
		}
	}

	// Example 4: Check minimum order quantity
	orderQty := types.MustDecimal("0.001")
	minQty := types.MustDecimal("0.01")

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
	defaultFee := types.ZeroDecimal
	if defaultFee.IsZero() {
		fmt.Println("\n✓ Using zero fee (maker rebate)")
	}

	// Example 6: Various zero representations
	fmt.Println("\n=== Zero Representations ===")
	zeroValues := []types.Decimal{
		types.MustDecimal("0"),
		types.MustDecimal("0.0"),
		types.MustDecimal("0.00"),
		types.MustDecimal("0.000000"),
		types.ZeroDecimal,
	}

	for _, val := range zeroValues {
		fmt.Printf("Decimal(%s).IsZero() = %v\n", val.String(), val.IsZero())
	}

	// Example 7: Non-zero values
	fmt.Println("\n=== Non-Zero Values ===")
	nonZeroValues := []types.Decimal{
		types.MustDecimal("0.0001"),
		types.MustDecimal("1"),
		types.MustDecimal("-0.01"),
		types.MustDecimal("100.50"),
	}

	for _, val := range nonZeroValues {
		fmt.Printf("Decimal(%s).IsZero() = %v\n", val.String(), val.IsZero())
	}

	// Example 8: After arithmetic operations
	fmt.Println("\n=== After Arithmetic ===")
	a := types.MustDecimal("10.5")
	b := types.MustDecimal("10.5")
	diff, _ := a.Sub(b)
	fmt.Printf("Decimal(%s) - Decimal(%s) = %s, IsZero=%v\n",
		a, b, diff, diff.IsZero())

	c := types.MustDecimal("0.1")
	d := types.MustDecimal("10")
	product, _ := c.Mul(d)
	one := types.MustDecimal("1")
	finalDiff, _ := product.Sub(one)
	fmt.Printf("Decimal(%s) × Decimal(%s) - Decimal(%s) = %s, IsZero=%v\n",
		c, d, one, finalDiff, finalDiff.IsZero())

	fmt.Println("\n=== Examples Complete ===")
}
