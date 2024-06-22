package middleware

import (
	"marketplace-system/config"
	customerror "marketplace-system/lib/customerrors"
	"marketplace-system/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Phone    string `json:"phone"`
	Id       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	jwt.StandardClaims
}

var (
	ErrNoAuthHeader      = customerror.NewInternalErrorf("Authorization header is missing")
	ErrInvalidAuthHeader = customerror.NewInternalErrorf("Authorization header is malformed")
	ErrClaimsInvalid     = customerror.NewInternalErrorf("Provided claims do not match expected scopes")
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

// GetJWSFromRequest extracts a JWS string from an Authorization: Bearer <jws> header
func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}
