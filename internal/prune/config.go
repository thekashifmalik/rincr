package prune

import (
	"github.com/thekashifmalik/rincr/internal/args"
)

type Config struct {
	Hourly  int
	Daily   int
	Monthly int
	Yearly  int
}

func NewConfig(args *args.Parsed) *Config {
	hourly, ok := args.GetOptionInt("hourly")
	if !ok {
		hourly = 24
	}
	daily, ok := args.GetOptionInt("daily")
	if !ok {
		daily = 30
	}
	monthly, ok := args.GetOptionInt("monthly")
	if !ok {
		monthly = 12
	}
	yearly, ok := args.GetOptionInt("yearly")
	if !ok {
		yearly = 10
	}
	return &Config{
		Hourly:  hourly,
		Daily:   daily,
		Monthly: monthly,
		Yearly:  yearly,
	}
}
