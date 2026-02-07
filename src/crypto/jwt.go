package crypto

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims specifies custom claims
type CustomClaims struct {
	SessionId string         `json:"session_id"`
	Role      string         `json:"role"`
	UserId    string         `json:"user_id"`
	CompanyId string         `json:"company_id"`
	Data      map[string]any `json:"data,omitempty"`
	jwt.RegisteredClaims
}

// Generate generates jwt token
func Generate(signingKey string, claims jwt.Claims) (string, error) {
	tn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := tn.SignedString([]byte(signingKey))
	return signedString, err
}

func JwtValid[T jwt.Claims](signedToken string, signingKey string, claims T) (T, bool, error) {
	token, err := jwt.ParseWithClaims(signedToken, claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signingKey), nil
		})

	var zero T
	if err != nil {
		return zero, false, err
	}

	if parsedClaims, ok := token.Claims.(T); ok && token.Valid {
		return parsedClaims, true, nil
	}

	return zero, false, err
}

// @Depricated use JwtValid instead
func Valid(signedToken string, signingKey string) (CustomClaims, bool, error) {
	token, err := jwt.ParseWithClaims(signedToken, &CustomClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signingKey), nil
		})
	if err != nil {
		return CustomClaims{}, false, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return *claims, true, nil
	}
	return CustomClaims{}, false, err
}
func SignAccessToken(AccessSecret string, expMin int, claims *CustomClaims) (string, error) {
	//claims.ExpiresAt= 20
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expMin)))
	token, err := Generate(AccessSecret, claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
func SignRefreshToken(RefreshSecret string, expMin int, claims *CustomClaims) (string, error) {
	//claims.ExpiresAt= 20
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expMin)))
	token, err := Generate(RefreshSecret, claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
