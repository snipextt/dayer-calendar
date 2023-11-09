package utils

import (
	"context"
	"time"
)

func NewContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	return ctx, cancel
}
