package middleware

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips   map[string]*rate.Limiter
	mu    *sync.RWMutex
	rate  rate.Limit
	burst int
}

// NewIPRateLimiter 创建一个新的 IP 限流器
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips:   make(map[string]*rate.Limiter),
		mu:    &sync.RWMutex{},
		rate:  r,
		burst: b,
	}
}

// GetLimiter 获取特定 IP 的限流器
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.ips[ip] = limiter
	}

	return limiter
}

func PathNormalizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		normalizedPath := strings.ReplaceAll(path, "//", "/")
		for strings.Contains(normalizedPath, "//") {
			normalizedPath = strings.ReplaceAll(normalizedPath, "//", "/")
		}
		c.Request.URL.Path = normalizedPath
		c.Next()
	}
}

// RateLimit 创建一个包含全局限流和 IP 限流的中间件
func RateLimit(globalLimiter *rate.Limiter, ipLimiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 检查全局限流
		if !globalLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}

		// 2. 检查 IP 限流
		ip := c.ClientIP()
		limiter := ipLimiter.GetLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests from your IP",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
