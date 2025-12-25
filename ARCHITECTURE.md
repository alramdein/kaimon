# Kaimon Architecture Overview

## Request Flow

```
Client Request
    │
    ▼
┌──────────────────────────────────────┐
│      Kaimon API Gateway              │
│                                      │
│  ┌────────────────────────────────┐ │
│  │   Framework Layer (Echo)       │ │
│  │   - HTTP Server                │ │
│  │   - Request/Response handling  │ │
│  └────────────────┬───────────────┘ │
│                   │                  │
│  ┌────────────────▼───────────────┐ │
│  │   Middleware Chain             │ │
│  │                                │ │
│  │  OnRequest Phase:              │ │
│  │   1. Logger    ────────────┐   │ │
│  │   2. CORS      ────────┐   │   │ │
│  │   3. Custom... ────┐   │   │   │ │
│  │                    │   │   │   │ │
│  │  Handler Execution │   │   │   │ │
│  │                    │   │   │   │ │
│  │  OnResponse Phase: │   │   │   │ │
│  │   1. Timer     ◄───┘   │   │   │ │
│  │   2. Custom... ◄───────┘   │   │ │
│  │                ◄───────────┘   │ │
│  └────────────────┬───────────────┘ │
│                   │                  │
│  ┌────────────────▼───────────────┐ │
│  │   Route Handler                │ │
│  │   - URL matching               │ │
│  │   - Proxy to backend           │ │
│  └────────────────┬───────────────┘ │
└───────────────────┼──────────────────┘
                    │
                    ▼
          Backend Services
          (users, products, etc.)
```

## Module Structure

### 1. pkg/framework
**Purpose**: Framework abstraction layer  
**Components**:
- `framework.go` - Interfaces (Framework, Router, Context, HandlerFunc)
- `echo.go` - Echo implementation

**Why**: Allows switching web frameworks without changing business logic.

### 2. pkg/routes
**Purpose**: Route configuration and compilation  
**Components**:
- `types.go` - Route data structures
- `compiler.go` - Compiles domain routes into single config
- `loader.go` - Loads compiled routes and registers with framework

**Flow**:
```
config/routes/users.json  ─┐
config/routes/products.json├─► Compiler ──► config/routes.json ──► Loader ──► Framework
config/routes/orders.json ─┘
```

### 3. pkg/middleware
**Purpose**: Middleware chain management  
**Components**:
- `types.go` - Middleware interfaces
- `manager.go` - Registration and execution order

**Phases**:
- **OnRequest**: Executed before handler (logging, auth, validation)
- **OnResponse**: Executed after handler (metrics, response transformation)

### 4. internal/middlewares
**Purpose**: Concrete middleware implementations  
**Structure**:
```
internal/middlewares/
├── onRequest/
│   ├── logger.go    # Request logging
│   ├── cors.go      # CORS headers
│   └── auth.go      # Authentication (example)
└── onResponse/
    ├── timer.go     # Response time tracking
    └── cache.go     # Response caching (example)
```

## Configuration System

### Domain Route Config
Each domain has its own JSON config:

```json
{
  "domain": "users",
  "basePath": "/api/v1/users",
  "globalLayers": ["logger", "cors"],
  "headers": {
    "X-Service": "users"
  },
  "routes": [...]
}
```

### Middleware Order Config
Controls execution order:

```json
{
  "onRequest": ["logger", "cors", "auth"],
  "onResponse": ["timer", "cache"]
}
```

Order matters! Middlewares execute in sequence.

## Extension Points

### Add New Framework
1. Implement `framework.Framework` interface
2. Implement `framework.Router` interface  
3. Implement `framework.Context` interface
4. Replace in `cmd/root.go`

### Add New Middleware
1. Create file in `internal/middlewares/onRequest/` or `onResponse/`
2. Implement `OnRequestMiddleware` or `OnResponseMiddleware` interface
3. Register in `cmd/root.go`
4. Add to execution order in `config/middleware-order.json`

### Add New Domain Routes
1. Create `config/routes/<domain>.json`
2. Run `./kaimon compile`
3. Restart gateway

## Why This Architecture?

✅ **Modular**: Each component has single responsibility  
✅ **Testable**: Interfaces allow easy mocking  
✅ **Extensible**: Add features without changing core  
✅ **Maintainable**: Clear separation of concerns  
✅ **Flexible**: Framework-agnostic design  
✅ **Simple**: Straightforward configuration
