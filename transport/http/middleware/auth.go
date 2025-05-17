package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/farganamar/evv-service/configs"
	"github.com/farganamar/evv-service/helpers/auth"
	"github.com/farganamar/evv-service/helpers/failure"
	"github.com/farganamar/evv-service/infras"
	"github.com/farganamar/evv-service/transport/http/response"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidToken    = errors.New("TOKEN_INVALID")
	ErrExpiredToken    = errors.New("TOKEN_EXPIRED")
	ErrUserNotVerified = errors.New("USER_NOT_VERIFIED")
	ErrUnauthorized    = errors.New("UNAUTHORIZED")
)

type AuthMiddlewareInterface interface {
	Authentication(p *ParamAuth) func(next http.Handler) http.Handler
}

type ParamAuth struct {
	IsVerified bool
	Roles      []string
	HasAccess  []string
}

type AuthValue struct {
	User       *auth.CustomClaims `json:"user"`
	Permission Permission         `json:"permission"`
}

type AuthMiddleware struct {
	tokenHelper *auth.TokenService
	config      *configs.Config
	redisConfig *infras.Redis
}

type Permission struct {
	Access      []string          `json:"access"`
	Entitlement map[string]string `json:"entitlement"`
	Usage       map[string]string `json:"usage"`
}

func NewAuthMiddleware(tokenHelper *auth.TokenService, cfg *configs.Config, redis *infras.Redis) *AuthMiddleware {
	return &AuthMiddleware{
		tokenHelper: tokenHelper,
		config:      cfg,
		redisConfig: redis,
	}
}

type ContextKeyType string

var ContextKey = ContextKeyType("user")

func (m *AuthMiddleware) Authentication(p *ParamAuth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
			if token == "" {
				response.WithError(w, failure.Unauthorized("Failed to get token"))
				return
			}

			claims, err := m.tokenHelper.ValidateToken(auth.ValidateToken{
				TokenString: token,
			})

			if err != nil {
				log.Error().Err(err).Msg("failed to validate token")
				response.WithJSONCodeStatus(w, http.StatusUnauthorized, nil, ErrInvalidToken.Error(), err.Error())
				return
			}

			if claims.IsVerified != p.IsVerified {
				response.WithJSONCodeStatus(w, http.StatusUnauthorized, nil, ErrInvalidToken.Error(), ErrUserNotVerified.Error())
				return
			}

			if len(p.Roles) > 0 {
				isRoleValid := false
				hasRole := make(map[string]bool)
				for _, role := range p.Roles {
					hasRole[role] = true
				}
				for _, r := range claims.Roles {
					if hasRole[r] {
						isRoleValid = true
						break
					}
				}

				if !isRoleValid {
					response.WithJSONCodeStatus(w, http.StatusUnauthorized, nil, ErrInvalidToken.Error(), ErrUnauthorized.Error())
					return
				}
			}

			if len(p.HasAccess) > 0 {
				newAccess := make(map[string]bool)
				for _, v := range p.HasAccess {
					newAccess[v] = true
				}
				isFeatureValid := false

				if !isFeatureValid {
					response.WithJSONCodeStatus(w, http.StatusUnauthorized, nil, ErrInvalidToken.Error(), ErrUnauthorized.Error())
					return
				}
			}

			user := AuthValue{
				User: claims,
			}

			ctx := context.WithValue(r.Context(), ContextKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
