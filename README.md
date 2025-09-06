# limitron

**limitron** is a Go rate limiting library supporting multiple algorithms and storage backends.  
Suitable for API protection, internal service throttling, or resource limiting with in-memory or Redis store.

## Features

- Rate limiting algorithms:
    - Token Bucket (`limitron.AlgorithmTokenBucket`)
    - Fixed Window (`limitron.AlgorithmFixedWindow`)
    - Sliding Window (`limitron.AlgorithmSlidingWindow`)
    - Leaky Bucket (`limitron.AlgorithmLeakyBucket`)
- Storage backends:
    - Memory (`limitron.StoreTypeMemory`)
    - Redis (`limitron.StoreTypeRedis`)
- Factory functions: create limiters with configuration, return `Limiter` interface
- Simple integration and usage pattern

## Installation

```sh
go get github.com/n0l3r/limitron
```

## Usage

### Configuration

Use constants for algorithm and store types to avoid typos.

```go
import "github.com/n0l3r/limitron"

cfg := limitron.Config{
    Algorithm:  limitron.AlgorithmTokenBucket, // or other algorithms
    Rate:       5,      // tokens per second / window size / leak rate
    Capacity:   10,     // bucket/window size
    StoreType:  limitron.StoreTypeMemory, // or limitron.StoreTypeRedis
    RedisClient: nil,   // set if using Redis
}
```

### Example: Token Bucket with Memory Store

```go
import (
    "github.com/n0l3r/limitron"
)

func main() {
    cfg := limitron.Config{
        Algorithm:  limitron.AlgorithmTokenBucket,
        Rate:       5,
        Capacity:   10,
        StoreType:  limitron.StoreTypeMemory,
    }

    limiter, err := limitron.NewLimiter(cfg)
    if err != nil {
        panic(err)
    }

    allowed, err := limiter.Allow("user-id")
    if err != nil {
        // handle error
    }
    if allowed {
        // proceed
    } else {
        // rate limited
    }
}
```

### Example: Token Bucket with Redis Store

```go
import (
    "github.com/go-redis/redis/v8"
    "github.com/n0l3r/limitron"
)

func main() {
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    cfg := limitron.Config{
        Algorithm:   limitron.AlgorithmTokenBucket,
        Rate:        5,
        Capacity:    10,
        StoreType:   limitron.StoreTypeRedis,
        RedisClient: rdb,
    }

    limiter, err := limitron.NewLimiter(cfg)
    if err != nil {
        panic(err)
    }
    allowed, err := limiter.Allow("user-id")
    // handle allowed
}
```

### Example: Leaky Bucket with Memory Store

```go
cfg := limitron.Config{
    Algorithm:  limitron.AlgorithmLeakyBucket,
    Rate:       1, // leak rate per second
    Capacity:   3,
    StoreType:  limitron.StoreTypeMemory,
}

limiter, err := limitron.NewLimiter(cfg)
allowed, _ := limiter.Allow("user-id")
```

### Using Concrete Limiter (advanced)

If you need algorithm-specific methods, use the dedicated factory:

```go
tb, err := limitron.NewTokenBucketLimiter(cfg)
allowed, _ := tb.Allow("user-id")
```

---

## Gin Middleware Example

You can easily use limitron as a middleware in [Gin](https://github.com/gin-gonic/gin) web framework to protect your endpoints.

```go
import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/n0l3r/limitron"
)

func RateLimitMiddleware(limiter limitron.Limiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Use c.ClientIP() or c.GetHeader("Authorization") or any unique key per user
        key := c.ClientIP() 
        allowed, err := limiter.Allow(key)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if !allowed {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
            return
        }
        c.Next()
    }
}

func main() {
    cfg := limitron.Config{
        Algorithm:  limitron.AlgorithmTokenBucket,
        Rate:       5,
        Capacity:   10,
        StoreType:  limitron.StoreTypeMemory,
    }
    limiter, err := limitron.NewLimiter(cfg)
    if err != nil {
        panic(err)
    }

    r := gin.Default()
    r.Use(RateLimitMiddleware(limiter))
    r.GET("/hello", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, World!")
    })
    r.Run(":8080")
}
```

**Tips:**
- Use `c.ClientIP()` for simple global rate limiting, or combine with other user identifiers for per-user rate limiting.
- Customize the middleware as needed for your application.

---

## API Overview

```go
type Limiter interface {
Allow(key string) (bool, error)
}

func NewLimiter(cfg Config) (Limiter, error)
func NewTokenBucketLimiter(cfg Config) (Limiter, error)
func NewFixedWindowLimiter(cfg Config) (Limiter, error)
func NewSlidingWindowLimiter(cfg Config) (Limiter, error)
func NewLeakyBucketLimiter(cfg Config) (Limiter, error)
```

### Constants

```go
const (
AlgorithmLeakyBucket   = "leaky_bucket"
AlgorithmTokenBucket   = "token_bucket"
AlgorithmFixedWindow   = "fixed_window"
AlgorithmSlidingWindow = "sliding_window"

StoreTypeMemory = "memory"
StoreTypeRedis  = "redis"
)
```

## Testing

Unit tests are provided for each algorithm in the `algorithms/` directory.

```sh
go test ./algorithms/...
```

## Contributing

Pull requests and discussions are welcome!  
Please report bugs or feature requests via Issues.

## License

MIT License