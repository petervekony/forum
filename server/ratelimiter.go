package server

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// IPRateLimiter is a struct that holds the connected IP addresses, a RWMutex, a rate.Limit on requests/second and a limit on package bursts
type IPRateLimiter struct {
	visitors         map[string]*Visitor
	lastVisitorCheck time.Time
	mu               *sync.RWMutex
	r                rate.Limit
	b                int
}

type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewIPRateLimiter returns a pointer to a new IPRateLimiter object
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		visitors: make(map[string]*Visitor),
		mu:       &sync.RWMutex{},
		r:        r,
		b:        b,
	}

	return i
}

// AddIP adds the IP adress to the connection map in IPRateLimiter
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()
	limiter := rate.NewLimiter(i.r, i.b)
	i.visitors[ip].limiter = limiter
	return limiter
}

// GetLimiter returns a pointer to the rate.Limiter in the IPRateLimiter connection map. If it doesn't exist, it adds it.
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	// Every time GetLimiter is called, we check if the connection map has been checked in the last 5 minutes. If not, it goes through all the connections,
	// If a connection wasn't used in the last 30 minutes, we delete the visitor from the connection map to avoid memory leaks.
	if time.Since(i.lastVisitorCheck) > 5*time.Minute {
		for ip, v := range i.visitors {
			if time.Since(v.lastSeen) > 30*time.Minute {
				delete(i.visitors, ip)
			}
		}
		i.lastVisitorCheck = time.Now()
	}
	limiter, exists := i.visitors[ip]
	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}
	i.mu.Unlock()
	i.visitors[ip].lastSeen = time.Now()
	return limiter.limiter
}

// LimitMiddleware is a middleware handler function that handles connections using a token bucket rate-limiter algorithm,
// if there are too many connections from the same IP address, the client gets a StatusTooManyRequests error.
func (i *IPRateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := i.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
