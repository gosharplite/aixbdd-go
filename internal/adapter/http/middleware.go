package http

import (
	"context"
	"net/http"
	"strings"

	"crm/internal/domain"
)

type contextKey string

const actorContextKey contextKey = "current_actor"

func AuthMiddleware(userRepo domain.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actor := &domain.CurrentActor{
				IsAuthenticated: false,
			}

			// 1. Check X-User-Name header (Direct testing / role simulation)
			userNameHeader := r.Header.Get("X-User-Name")
			roleHeader := r.Header.Get("X-User-Role")

			// 2. Check Authorization Bearer header
			authHeader := r.Header.Get("Authorization")
			if userNameHeader == "" && strings.HasPrefix(authHeader, "Bearer ") {
				userNameHeader = strings.TrimPrefix(authHeader, "Bearer ")
			}

			if userNameHeader != "" {
				user, err := userRepo.GetByName(r.Context(), userNameHeader)
				if err == nil && user != nil {
					actor.User = user
					actor.IsAuthenticated = true
					// Allow role override if specified in header
					if roleHeader != "" {
						actor.User.Role = domain.UserRole(roleHeader)
					}
				} else {
					// Auto-create or temporary actor for simulation
					role := domain.RoleSalesRep
					if roleHeader != "" {
						role = domain.UserRole(roleHeader)
					}
					actor.User = &domain.User{
						ID:       99,
						Username: strings.ToLower(userNameHeader),
						Name:     userNameHeader,
						Role:     role,
					}
					actor.IsAuthenticated = true
				}
			}

			ctx := context.WithValue(r.Context(), actorContextKey, actor)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actor := GetActor(r)
		if actor == nil || !actor.IsAuthenticated {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "需要登入以取得權限")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetActor(r *http.Request) *domain.CurrentActor {
	actor, ok := r.Context().Value(actorContextKey).(*domain.CurrentActor)
	if !ok || actor == nil {
		return &domain.CurrentActor{IsAuthenticated: false}
	}
	return actor
}
