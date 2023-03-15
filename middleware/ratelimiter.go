package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/**
 * Created by @CaomaoBoy on 2023/3/14.
 *  email:<115882934@qq.com>
 */

type RateLimiter struct {
	MaxRequests int           // 最大请求数
	WindowSize  time.Duration // 时间窗口大小
	IPMap       map[string]int64   // 存储每个IP地址对应的请求次数和最近一次请求时间
}
func (rl *RateLimiter) LimitHandler(c *gin.Context) {
	ip := c.ClientIP()
	now := time.Now().UnixNano()
	if _, ok := rl.IPMap[ip]; !ok {
		rl.IPMap[ip] = now
		return
	}
	count, lastTime := rl.IPMap[ip], rl.WindowSize.Nanoseconds()+rl.IPMap[ip]

	if now < lastTime {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	count++
	if int(count) > rl.MaxRequests {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	rl.IPMap[ip] = now
}