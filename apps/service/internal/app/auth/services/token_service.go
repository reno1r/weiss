package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/reno1r/weiss/apps/service/internal/config"
)

type TokenService struct {
	secretKey        []byte
	issuer           string
	accessExpiresIn  time.Duration
	refreshExpiresIn time.Duration
}

func NewTokenService(config *config.Config) (*TokenService, error) {
	if config.JwtSecret == "" {
		return nil, errors.New("JWT secret is required")
	}

	if config.JwtIssuer == "" {
		return nil, errors.New("JWT issuer is required")
	}

	accessExpiresIn, err := time.ParseDuration(config.JwtAccessExpires)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT access expiration: %w", err)
	}

	refreshExpiresIn, err := time.ParseDuration(config.JwtRefreshExpires)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT refresh expiration: %w", err)
	}

	return &TokenService{
		secretKey:        []byte(config.JwtSecret),
		issuer:           config.JwtIssuer,
		accessExpiresIn:  accessExpiresIn,
		refreshExpiresIn: refreshExpiresIn,
	}, nil
}

type Claims struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Type  string `json:"type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (ts *TokenService) GenerateAccessToken(userID uint64, email, phone string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(ts.accessExpiresIn)
	jti, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT ID: %w", err)
	}

	claims := &Claims{
		Email: email,
		Phone: phone,
		Type:  "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ts.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        jti.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ts.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ts *TokenService) GenerateRefreshToken(userID uint64, email, phone string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(ts.refreshExpiresIn)
	jti, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT ID: %w", err)
	}

	claims := &Claims{
		Email: email,
		Phone: phone,
		Type:  "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ts.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        jti.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ts.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ts *TokenService) GenerateTokenPair(userID uint64, email, phone string) (*TokenPair, error) {
	accessToken, err := ts.GenerateAccessToken(userID, email, phone)
	if err != nil {
		return nil, err
	}

	refreshToken, err := ts.GenerateRefreshToken(userID, email, phone)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ts *TokenService) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return ts.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (ts *TokenService) VerifyRefreshToken(tokenString string) (*Claims, error) {
	claims, err := ts.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func (ts *TokenService) GetUserID(claims *Claims) (uint64, error) {
	if claims.Subject == "" {
		return 0, errors.New("subject claim is missing")
	}

	var userID uint64
	_, err := fmt.Sscanf(claims.Subject, "%d", &userID)
	if err != nil {
		return 0, fmt.Errorf("invalid subject claim: %w", err)
	}

	return userID, nil
}

func (ts *TokenService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	claims, err := ts.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID, err := ts.GetUserID(claims)
	if err != nil {
		return nil, err
	}

	return ts.GenerateTokenPair(userID, claims.Email, claims.Phone)
}
