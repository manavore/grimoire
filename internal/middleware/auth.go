package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/manavore/grimoire/internal/services/auth"
)

type UserContextKey string

const (
	UserInfoKey UserContextKey = "userinfo"
)

type AuthMidlleware struct {
	AuthService *auth.AuthService
}

func NewAuthMiddleware(authService *auth.AuthService) *AuthMidlleware {
	return &AuthMidlleware{
		AuthService: authService,
	}
}

func (m *AuthMidlleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			authHeader = "Bearer " + cookie.Value
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userInfo, err := m.AuthService.GetUserInfo(tokenString)
		if err != nil {
			// Invalid or expired token
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserInfoKey, userInfo)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func GetUserInfo(ctx context.Context) (*auth.UserInfo, bool) {
	userInfo, ok := ctx.Value(UserInfoKey).(*auth.UserInfo)

	return userInfo, ok
}
