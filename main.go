package main

import (
	"context"

	"github.com/ihezebin/oneness/logger"
	"github.com/ihezebin/promisys/cmd"
)

func main() {
	ctx := context.Background()
	if err := cmd.Run(ctx); err != nil {
		logger.Fatalf(ctx, "cmd run error: %v", err)
	}
}
