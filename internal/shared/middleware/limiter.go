package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type visitor struct {
	count   int
	resetAt time.Time
}

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	max      int
	window   time.Duration
}

func NewLimiterMiddleware(max int, window time.Duration) func(huma.Context, func(huma.Context)) {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		max:      max,
		window:   window,
	}

	go rl.cleanup()

	return func(ctx huma.Context, next func(huma.Context)) {
		ip, _, err := net.SplitHostPort(ctx.RemoteAddr())
		if err != nil {
			ip = ctx.RemoteAddr()
		}

		rl.mu.Lock()
		v, ok := rl.visitors[ip]
		now := time.Now()
		if !ok || now.After(v.resetAt) {
			v = &visitor{count: 0, resetAt: now.Add(rl.window)}
			rl.visitors[ip] = v
		}
		v.count++
		count := v.count
		resetAt := v.resetAt
		rl.mu.Unlock()

		remaining := rl.max - count
		if remaining < 0 {
			remaining = 0
		}

		ctx.SetHeader("X-RateLimit-Limit", strconv.Itoa(rl.max))
		ctx.SetHeader("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if count > rl.max {
			retryAfter := int(time.Until(resetAt).Seconds()) + 1
			ctx.SetHeader("Retry-After", strconv.Itoa(retryAfter))
			ctx.SetHeader("Content-Type", "application/json")
			ctx.SetStatus(http.StatusTooManyRequests)
			_, _ = ctx.BodyWriter().Write([]byte(`{"title":"Too Many Requests","status":429}`))
			return
		}

		next(ctx)
	}
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.After(v.resetAt) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
