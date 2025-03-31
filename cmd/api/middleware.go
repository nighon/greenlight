package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		// Maps are not thread-safe. We need to lock the mutex before reading from the map.
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Background goroutine to clean up old clients.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.config.limiter.enabled {
			next.ServeHTTP(w, r)
			return
		}

		ip := realip.FromRequest(r)

		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
			}
		}

		clients[ip].lastSeen = time.Now()

		// Check if the request is allowed by the rate limiter.
		// `limiter.Allow()` returns `true` if the request is allowed, and `false` if the request is not allowed.
		// It consumes a token from the bucket.
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		// We aren't using defer here because we want to unlock the mutex as soon as possible.
		// If we deferred the unlock, it only would be executed after all the downstream handlers have returned.
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
