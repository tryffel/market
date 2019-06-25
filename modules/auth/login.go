package auth

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/market/config"
	"github.com/tryffel/market/modules/Error"
	"github.com/tryffel/market/modules/request"
	"github.com/tryffel/market/modules/response"
	"github.com/tryffel/market/modules/util"
	"github.com/tryffel/market/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	ClaimUserId   string = "user_id"
	ClaimUsername string = "user_name"
	ClaimNonce    string = "nonce"
)

type LoginForm struct {
	Username string
	Password string
}

type TokenForm struct {
	Token string
}

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

type Auth struct {
	store          *storage.Store
	pKey           string
	expireTokens   bool
	expiryDuration time.Duration
}

func NewAuth(store *storage.Store, conf *config.Config) Auth {
	a := Auth{
		store:          store,
		pKey:           conf.Tokens.Key,
		expireTokens:   conf.Tokens.Expire,
		expiryDuration: conf.Tokens.Interval.ToDuration(),
	}
	return a
}

func (a *Auth) Login(resp response.Response, req request.Request) {
	headers := req.Headers()
	token := headers[config.AuthorizationHeader]
	if len(token) > 0 {
		print("User already has token!")
		err := response.BadRequest("Already logged in", resp)
		if err != nil {
			logrus.Error(err)
		}
	}

	dto := LoginForm{}
	err := json.NewDecoder(req.Body()).Decode(&dto)
	if err != nil {
		logrus.Error(err)
	}

	user, err := a.store.User.FindByName(dto.Username)
	if err == nil && user != nil {
		if PasswordMatches(dto.Password, user.Password) {
			token, err := NewToken(user.Id, a.pKey,
				util.RandomKey(20), a.expiryDuration)
			if err != nil {
				logrus.Error(err)
			}

			t := TokenForm{
				Token: token,
			}

			response.JsonOk(t, resp)

		} else {
			response.Unauthorized(resp)

		}
	} else {
		response.Unauthorized(resp)
	}
}
