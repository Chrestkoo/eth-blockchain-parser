package middleware

import (
	"net/http"
	"sync"
	"time"
)

var ipRequestMap sync.Map

// Rate limit configuration
const (
	maxRequests     = 100              // Maximum allowed requests
	timeWindow      = 60 * time.Second // Time window for rate limiting
	cleanupInterval = 5 * time.Minute  // Interval to clean up stale IPs
)

// ClientInfo stores the number of requests and the last access time for each IP
type ClientInfo struct {
	Requests    int
	LastRequest time.Time
}

// Middleware to block requests based on IP rate limiting
func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		// Check if the IP has made requests before
		if clientInfo, ok := ipRequestMap.Load(ip); ok {
			info := clientInfo.(*ClientInfo)

			// Check if the client is within the time window
			if time.Since(info.LastRequest) < timeWindow {
				// If requests exceed the limit, block the IP
				if info.Requests >= maxRequests {
					http.Error(w, "Too many requests", http.StatusTooManyRequests)
					return
				}
				info.Requests++
			} else {
				// Reset the count after the time window has passed
				info.Requests = 1
			}
			info.LastRequest = time.Now()
		} else {
			// First request from this IP
			ipRequestMap.Store(ip, &ClientInfo{
				Requests:    1,
				LastRequest: time.Now(),
			})
		}

		// Pass to the next handler if within limit
		next.ServeHTTP(w, r)
	})
}

// Periodically clean up stale IP entries to prevent memory leaks
func CleanupIPs() {
	for {
		time.Sleep(cleanupInterval)
		ipRequestMap.Range(func(key, value interface{}) bool {
			clientInfo := value.(*ClientInfo)
			// Remove entries older than the time window
			if time.Since(clientInfo.LastRequest) > timeWindow {
				ipRequestMap.Delete(key)
			}
			return true
		})
	}
}
