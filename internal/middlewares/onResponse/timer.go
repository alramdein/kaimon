package onresponse

import (
	"log"
	"time"

	"github.com/alramdein/kaimon/pkg/framework"
	"github.com/alramdein/kaimon/pkg/middleware"
)

func init() {
	// Self-register this middleware
	middleware.RegisterOnResponseFactory("timer", func() middleware.OnResponseMiddleware {
		return NewTimerMiddleware()
	})
}

// TimerMiddleware measures request duration
type TimerMiddleware struct{}

func NewTimerMiddleware() *TimerMiddleware {
	return &TimerMiddleware{}
}

func (m *TimerMiddleware) Name() string {
	return "timer"
}

func (m *TimerMiddleware) Handle(next framework.HandlerFunc) framework.HandlerFunc {
	return func(ctx framework.Context) error {
		start := time.Now()
		err := next(ctx)
		duration := time.Since(start)

		req := ctx.Request()
		log.Printf("[TIMER] %s %s took %v", req.Method, req.URL.Path, duration)

		return err
	}
}
