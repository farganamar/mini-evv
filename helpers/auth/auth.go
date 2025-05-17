package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/farganamar/evv-service/configs"
	"github.com/farganamar/evv-service/infras"
)

type TokenServiceInterface interface {
	GenerateTokenPair(ctx context.Context, payload *ParamsGenerateToken) (*TokenPair, error)
	ValidateToken(paramValidate ValidateToken) (*CustomClaims, error)
	StoreToken(ctx context.Context, token string, tokenType string, ttl int) error
}

var (
	ErrInvalidToken = errors.New("TOKEN_INVALID")
	ErrExpiredToken = errors.New("TOKEN_EXPIRED")
)

type ParamsGenerateToken struct {
	UserID     string   `json:"user_id"`
	Email      string   `json:"email"`
	Username   string   `json:"username"`
	IsVerified bool     `json:"is_verified"`
	Roles      []string `json:"roles"`
}

type TokenService struct {
	redisClient *infras.Redis
	config      *configs.Config
}
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	IssuedAt     time.Time
}

type CustomClaims struct {
	UserID     string   `json:"user_id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	IsVerified bool     `json:"is_verified"`
	Roles      []string `json:"roles"`
	Type       string   `json:"type"`
	jwt.RegisteredClaims
}

type ValidateToken struct {
	TokenString string
	Type        string
	Secret      string
}

func NewTokenService(redis *infras.Redis, cfg *configs.Config) *TokenService {
	return &TokenService{
		redisClient: redis,
		config:      cfg,
	}
}

func (s *TokenService) GenerateTokenPair(ctx context.Context, payload *ParamsGenerateToken) (*TokenPair, error) {
	accessToken, err := s.createToken(payload, []byte(s.config.AccessToken.Secret), s.config.AccessToken.ExpiryInHour, "ACCESS_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.createToken(payload, []byte(s.config.RefreshToken.Secret), s.config.RefreshToken.ExpiryInHour, "REFRESH_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * time.Duration(s.config.RefreshToken.ExpiryInHour)),
		IssuedAt:     time.Now(),
	}, nil
}

func (s *TokenService) createToken(payload *ParamsGenerateToken, secret []byte, expiry int, tokenType string) (string, error) {
	claims := CustomClaims{
		UserID:     payload.UserID,
		Email:      payload.Email,
		Username:   payload.Username,
		IsVerified: true,
		Roles:      payload.Roles,
		Type:       tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiry))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (s *TokenService) StoreToken(ctx context.Context, token string, tokenType string, ttl int) error {
	var key string
	if ttl == 0 {
		ttl = s.config.RefreshToken.ExpiryInHour
	}
	switch tokenType {
	case "ACCESS_TOKEN":
		key = fmt.Sprintf("access_token:%s", token)
	case "REFRESH_TOKEN":
		key = fmt.Sprintf("refresh_token:%s", token)
	}

	return s.redisClient.RedisClient.Set(ctx, key, true, time.Hour*time.Duration(ttl)).Err()
}

func (s *TokenService) ValidateToken(paramValidate ValidateToken) (*CustomClaims, error) {
	var secret string
	var key string
	token, err := jwt.Parse(paramValidate.TokenString, func(token *jwt.Token) (interface{}, error) {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			tokenType, ok := claims["type"].(string)
			if !ok {
				return nil, ErrInvalidToken
			}

			if tokenType == "ACCESS_TOKEN" {
				secret = s.config.AccessToken.Secret
				key = fmt.Sprintf("access_token:%s", paramValidate.TokenString)
			} else if tokenType == "REFRESH_TOKEN" {
				secret = s.config.RefreshToken.Secret
				key = fmt.Sprintf("refresh_token:%s", paramValidate.TokenString)
			} else {
				return nil, ErrInvalidToken
			}
		}
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error().Msgf("Unexpected signing method: %v", token.Header["alg"])
			return nil, ErrInvalidToken
		}

		return []byte(secret), nil
	})

	if err != nil && jwt.ErrTokenExpired.Error() == "" {
		log.Error().Err(err).Msg("failed to parse token")
		return nil, err
	}

	if err != nil && jwt.ErrTokenExpired.Error() != "" {
		return nil, ErrExpiredToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	roles := claims["roles"].([]interface{})
	var rolesString []string
	for _, role := range roles {
		rolesString = append(rolesString, role.(string))
	}

	user := CustomClaims{
		UserID:     claims["user_id"].(string),
		Email:      claims["email"].(string),
		Username:   claims["username"].(string),
		IsVerified: claims["is_verified"].(bool),
		Roles:      rolesString,
		Type:       claims["type"].(string),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(int64(claims["exp"].(float64)), 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(int64(claims["iat"].(float64)), 0)),
		},
	}

	if user.RegisteredClaims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	if exist, err := s.redisClient.RedisClient.Exists(context.TODO(), key).Result(); err != nil {
		log.Error().Err(err).Msg("failed to get token from redis")
	} else if exist == 1 {
		return nil, ErrInvalidToken
	}

	return &user, nil
}
