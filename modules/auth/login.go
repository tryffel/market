package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/tryffel/market/modules/Error"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	ClaimUserId   string = "user_id"
	ClaimUsername string = "user_name"
	ClaimNonce    string = "nonce"
)

// Get password hash
func GetPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// PasswordMatches true|false
func PasswordMatches(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// NewToken issues a new token for user_id
// if ExpireDuration == 0, disable expiration
func NewToken(userId string, privateKey string, nonce string, expireDuration time.Duration) (string, error) {
	var token *jwt.Token = nil
	if expireDuration == 0 {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			ClaimUserId: userId,
			ClaimNonce:  nonce,
		})
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			ClaimUserId: userId,
			ClaimNonce:  nonce,
			"nbf":       time.Now(),
			"exp":       time.Now().Add(expireDuration).Unix(),
		})
	}

	tokenString, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

// ValidateToken validates and parses user id. If valid, return user_id, else return error description
// Return user_id, nonce, error
func ValidateToken(tokenString string, privateKey string,
	VerifyExpiresAt bool, expirationDuration time.Duration) (string, string, error) {

	// Fixme If token has expired-field, it is invalid regardless of enabling / disabling during validation

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(privateKey), nil
	})

	if err != nil {
		e, ok := err.(*jwt.ValidationError)
		if ok {
			if e.Inner.Error() == "Token is expired" {
				return "", "", Error.NewEndUserError(Error.Einvalid, "token expired")
			}

		} else {
			return "", "", Error.NewEndUserError(Error.Einvalid, "invalid token")
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		t := time.Now()
		if claims.VerifyExpiresAt(t.Unix(), VerifyExpiresAt) {
			user := claims[ClaimUserId].(string)
			nonce := claims[ClaimNonce].(string)
			return user, nonce, nil
		}
		return "", "", Error.NewEndUserError(Error.Einvalid, "token expired")
	}
	return "", "", Error.NewEndUserError(Error.Einvalid, "invalid token")
}
