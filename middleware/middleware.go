package middleware

import (
	"errors"
	"marketplace-system/config"
	utils "marketplace-system/lib/helper"
	"marketplace-system/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	Phone    string `json:"phone"`
	Id       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	jwt.StandardClaims
}

var (
	prefixBearer = "Bearer"
	headerAuth   = "Authorization"
)

var (
	ErrNoAuthHeader      = errors.New("Authorization header required")
	ErrInvalidAuthHeader = errors.New("Invalid authorization header")
	ErrInvalidToken      = errors.New("Invalid token")
	ErrMissingAuth       = errors.New("Missing Authorization header")
)

func GenerateToken(customer models.Customer, cfg config.Config) (string, error) {

	claims := &Claims{
		Phone:    customer.Phone,
		Email:    customer.Email,
		Id:       int(customer.CustomerID),
		FullName: customer.FullName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Hour).Unix(),
		},
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.SecretKeyJWT))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func JWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(headerAuth)
			if authHeader == "" {
				return utils.RespondWithError(c, http.StatusUnauthorized, utils.GetErrorResponse(ErrMissingAuth.Error(), http.StatusUnauthorized))
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != prefixBearer {
				return utils.RespondWithError(c, http.StatusUnauthorized, utils.GetErrorResponse(ErrInvalidToken.Error(), http.StatusUnauthorized))
			}

			tokenString := tokenParts[1]
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil {
				return utils.RespondWithError(c, http.StatusUnauthorized, utils.GetErrorResponse(ErrInvalidToken.Error(), http.StatusUnauthorized))
			}

			claims, ok := token.Claims.(*Claims)
			if !ok || !token.Valid {
				return utils.RespondWithError(c, http.StatusUnauthorized, utils.GetErrorResponse(ErrInvalidToken.Error(), http.StatusUnauthorized))
			}

			// Set user context
			c.Set("user", claims)
			c.Set("id", claims.Id)

			return next(c)
		}
	}
}
