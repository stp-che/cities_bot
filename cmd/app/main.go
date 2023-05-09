package main

import (
	"context"

	"github.com/stp-che/cities_bot/pkg/log"
)

func main() {
	ctx := context.Background()
	log.Info(ctx, "Cities Bot started")
	log.Info(ctx, "Cities Bot finished")
}
