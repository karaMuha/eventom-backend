package middlewares

import (
	"context"
	"eventom-backend/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
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
