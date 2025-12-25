# Kaimon API Gateway

A lean, simple, modular API gateway with framework abstraction built in Go.

## Features

- **Framework Agnostic**: Uses Echo by default, but easily swappable with any web framework
- **Modular Architecture**: Clean separation of concerns with pkg structure
- **JSON-based Routes**: Configure routes via JSON files per domain
- **Middleware Chain**: OnRequest and OnResponse middleware layers with configurable execution order
- **Route Compilation**: Compile multiple domain route configs into single optimized file

## Architecture

```
kaimon/
├── cmd/               # CLI commands (compile, serve)
├── config/
│   ├── routes/       # Domain-specific route configs
│   │   ├── users.json
│   │   └── products.json
│   ├── middleware-order.json  # Middleware execution order
│   └── routes.json   # Compiled routes (generated)
├── internal/
│   └── middlewares/
│       ├── onRequest/   # Request-phase middlewares
│       └── onResponse/  # Response-phase middlewares
└── pkg/
    ├── framework/    # Framework abstraction layer
    ├── middleware/   # Middleware manager
    └── routes/       # Route compiler & loader
```

## Getting Started

### 1. Build

```bash
go build -o kaimon
```

### 2. Configure Routes

Create route configs per domain in `config/routes/`:

**config/routes/users.json**:
```json
{
  "domain": "users",
  "basePath": "/api/v1/users",
  "globalLayers": ["logger", "cors"],
  "routes": [
    {
      "path": "",
      "method": "GET",
      "target": "http://localhost:8081/users"
    },
    {
      "path": "/:id",
      "method": "GET",
      "target": "http://localhost:8081/users/:id"
    }
  ]
}
```

### 3. Configure Middleware Order

**config/middleware-order.json**:
```json
{
  "onRequest": ["logger", "cors"],
  "onResponse": ["timer"]
}
```

### 4. Compile Routes

```bash
./kaimon compile
```

This compiles all domain routes into `config/routes.json`.

### 5. Start Gateway

```bash
./kaimon serve
```

Gateway starts on port 8080.

## Adding Custom Middlewares

### OnRequest Middleware

Create a file in `internal/middlewares/onRequest/`:

```go
package onrequest

import "github.com/alramdein/kaimon/pkg/framework"

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

Register in `cmd/root.go`:
```go
mwManager.RegisterOnRequest(onrequest.NewMyMiddleware())
```

### OnResponse Middleware

Same pattern, but in `internal/middlewares/onResponse/`.

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
