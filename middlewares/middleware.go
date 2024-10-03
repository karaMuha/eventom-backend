package middlewares

import (
	"context"
	"eventom-backend/utils"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

type Middleware func(http.Handler, *utils.Logger) http.Handler

func CreateStack(mws ...Middleware) Middleware {
	return func(next http.Handler, logger *utils.Logger) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			mw := mws[i]
			next = mw(next, logger)
		}

		return next
	}
}

func AuthMiddleware(next http.Handler, logger *utils.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTarget := r.Method + " " + strings.Split(r.URL.Path, "/")[1]

		if !utils.ProtectedRoutes[requestTarget] {
			next.ServeHTTP(w, r)
			return
		}

		jwtToken, err := r.Cookie("jwt")

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		verifiedToken, err := utils.VerifyJwt(jwtToken.Value)

		if err != nil {
			logger.Log(utils.LevelError, fmt.Sprintf("Failed to verify jwt: %s", err.Error()), map[string]string{
				"Request IP Address: ": r.RemoteAddr,
			})
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// extract user id from token for further usage
		claims, ok := verifiedToken.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "Could not convert jwt claims", http.StatusInternalServerError)
			return
		}

		userId, ok := claims["user_id"].(string)

		if !ok {
			http.Error(w, "Could not convert user id from jwt claims to string", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextUserIdKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RateLimiterMiddleware(next http.Handler, logger *utils.Logger) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var mutex sync.Mutex
	clients := make(map[string]*client)

	// loop through the client ips every minute and clean up those who haven't send a request for atleast three minutes
	go func() {
		for {
			time.Sleep(time.Minute)
			mutex.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logger.Log(utils.LevelError, err.Error(), map[string]string{
				"Request IP Address: ": ip,
				"Request URL: ":        r.URL.Path,
			})
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mutex.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(3, 30),
			}
		}

		clients[ip].lastSeen = time.Now()

		mutex.Unlock()

		if !clients[ip].limiter.Allow() {
			logger.Log(utils.LevelError, "Too many requests", map[string]string{
				"Request IP Address: ": ip,
				"Request URL: ":        r.URL.Path,
			})
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
