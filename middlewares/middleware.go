package middlewares

import (
	"eventom-backend/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		claims, _ := verifiedToken.Claims.(jwt.MapClaims)
		userId := claims["userId"].(string)
		r.Header.Set("userId", userId)
		next.ServeHTTP(w, r)
	})
}
