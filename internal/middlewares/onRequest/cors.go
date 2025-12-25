package onrequest

import (
	"github.com/alramdein/kaimon/pkg/framework"
	"github.com/alramdein/kaimon/pkg/middleware"
)

func init() {
	// Self-register this middleware
	middleware.RegisterOnRequestFactory("cors", func() middleware.OnRequestMiddleware {
		return NewCORSMiddleware()
	})
}

// CORSMiddleware handles CORS headers
type CORSMiddleware struct{}

func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

func (m *CORSMiddleware) Name() string {
	return "cors"
}

func (m *CORSMiddleware) Handle(next framework.HandlerFunc) framework.HandlerFunc {
	return func(ctx framework.Context) error {
		res := ctx.Response()
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if ctx.Request().Method == "OPTIONS" {
			res.WriteHeader(204)
			return nil
		}

		return next(ctx)
	}
}
