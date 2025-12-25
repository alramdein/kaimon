# Kaimon API Gateway

A lean, simple, modular API gateway with framework abstraction built in Go.

## Background

Today's API gateway landscape has two main players: [KrakenD](https://www.krakend.io/) and [Kong](https://konghq.com/). Both have their pros and cons.

**KrakenD**: Focuses on faster performance and minimal setup  
**Kong**: Provides comprehensive API management tools

**The Problem**:
- KrakenD: The `krakend.json` can be hell to manage
- Kong: Rich features can make the gateway bloated, adding latency and maintenance overhead

**My Solution**:

Kaimon sits exactly between those tools. I made it as simple as possible, while keeping the `routes.json` management clean and maintainable.

I also ditched the plugin architecture and put middlewares directly as a chain in the code. This makes things faster. I also made the API gateway agnostic of any framework, in whcih users can choose which web framework they want.

## Features

- **Framework Agnostic**: Uses Echo by default, but easily swappable with any web framework
- **Modular Architecture**: Clean separation of concerns with pkg structure
- **JSON-based Routes**: Configure routes via JSON files per domain
- **Middleware Chain**: OnRequest and OnResponse middleware layers with configurable execution order
- **Route Compilation**: Compile multiple domain route configs into single optimized file
- **No Plugins**: Middlewares run directly in-process for better performance

## Getting Started

### 1. Build

```bash
go build -o kaimon
```

### 2. Configure Global Settings

**config/global.json**:
```json
{
  "middlewares": {
    "onRequest": ["logger", "cors"],
    "onResponse": ["timer"]
  },
  "headers": {
    "X-Powered-By": "Kaimon"
  }
}
```

### 3. Configure Routes

Create route configs per domain in `config/routes/`:

**config/routes/users.json**:
```json
{
  "domain": "users",
  "basePath": "/api/v1/users",
  "middlewares": {
    "onRequest": [],
    "onResponse": []
  },
  "routes": [
    {
      "path": "",
      "method": "GET",
      "target": "http://localhost:8081/users"
    },
    {
      "path": "/:id",
      "method": "GET",
      "target": "http://localhost:8081/users/:id",
      "middlewares": {
        "onRequest": ["auth"],
        "onResponse": []
      }
    }
  ]
}
```

**Note**: Middlewares can be defined at three levels:
- **Global** (applies to all routes)
- **Domain** (applies to all routes in a domain)
- **Route** (applies to specific route)

### 4. Compile Routes

```bash
./kaimon compile
```

This compiles all domain routes and global config into `config/routes.json`.

### 5. Start Gateway

```bash
./kaimon serve
```

Gateway starts on port 8080.

## Adding Custom Middlewares

Middlewares auto-register themselves - just create the file and rebuild. No need to manually register anywhere!

### OnRequest Middleware

Create a file in `internal/middlewares/onRequest/`:

```go
package onrequest

import (
    "github.com/alramdein/kaimon/pkg/framework"
    "github.com/alramdein/kaimon/pkg/middleware"
)

func init() {
    // Auto-register on startup
    middleware.RegisterOnRequestFactory("my-middleware", func() middleware.OnRequestMiddleware {
        return NewMyMiddleware()
    })
}

type MyMiddleware struct{}

func NewMyMiddleware() *MyMiddleware {
    return &MyMiddleware{}
}

func (m *MyMiddleware) Name() string {
    return "my-middleware"
}

func (m *MyMiddleware) Handle(next framework.HandlerFunc) framework.HandlerFunc {
    return func(ctx framework.Context) error {
        // Pre-processing
        err := next(ctx)
        // Post-processing
        return err
    }
}
```

### OnResponse Middleware

Same pattern, but in `internal/middlewares/onResponse/` and use `RegisterOnResponseFactory`:

```go
func init() {
    middleware.RegisterOnResponseFactory("my-response-mw", func() middleware.OnResponseMiddleware {
        return NewMyResponseMiddleware()
    })
}
```

**That's it!** The middleware is now available to use in your route configs.

## Switching Frameworks

To use a different framework, implement the `framework.Framework` interface:

```go
type Framework interface {
    Router() Router
    Start(address string) error
    Shutdown() error
}
```

Then replace in `cmd/root.go`:
```go
// Instead of:
fw := framework.NewEchoFramework()

// Use:
fw := framework.NewMyFramework()
```

## Commands

- `./kaimon compile` - Compile route configurations
- `./kaimon serve` - Start API gateway server

## License

MIT
