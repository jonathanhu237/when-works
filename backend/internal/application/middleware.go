package application

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const requesterContextKey = contextKey("requester")

type RequesterInfo struct {
	UserID   uuid.UUID
	Username string
	IsAdmin  bool
}

// requireAuth middleware validates JWT token and sets user in context
func (app *Application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		cookie, err := r.Cookie("accessToken")
		if err != nil {
			app.unauthorizedResponse(w, r)
			return
		}

		// Parse and validate token
		tokenString := cookie.Value
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(app.config.JWT.Secret), nil
		})

		if err != nil {
			app.unauthorizedResponse(w, r)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			app.unauthorizedResponse(w, r)
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			app.unauthorizedResponse(w, r)
			return
		}

		// Set requester in context
		requester := &RequesterInfo{
			UserID:   userID,
			Username: claims.Username,
			IsAdmin:  claims.IsAdmin,
		}

		ctx := context.WithValue(r.Context(), requesterContextKey, requester)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireAdmin middleware ensures the user is an admin
func (app *Application) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requester := r.Context().Value(requesterContextKey).(*RequesterInfo)

		if !requester.IsAdmin {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
