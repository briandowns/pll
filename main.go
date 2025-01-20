package main

import (
	"context"
	"fmt"
	"os"

	"github.com/briandowns/pll/pll"
)

func main() {
	ctx := context.Background()

	token := os.Getenv("PLL_BEARER_TOKEN")

	p := pll.NewPLL(token)

	// standings, err := pll.Standings(ctx, 2023, false)
	// if err != nil {
	// 	fmt.Fprint(os.Stderr, err)
	// 	os.Exit(1)
	// }

	stats, err := p.PlayerStats(ctx, 2024, 5, "regular", pll.PlayerStatistics)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%#+v\n", stats)

	os.Exit(0)
}
