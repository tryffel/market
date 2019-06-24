package middleware

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tryffel/market/config"
	"github.com/tryffel/market/modules/Error"
	"github.com/tryffel/market/modules/response"
	"github.com/tryffel/market/storage"
	"net/http"
	"strings"
)

type Auth struct {
	store *storage.Store
	pKey  string
}

// NewAuth creates new auth instance
func NewAuth(config *config.Config, store *storage.Store) (Auth, error) {
	auth := &Auth{
		store: store,
	}

	if config.Tokens.Key == "" {
		return *auth, errors.New("server private key cannot be empty")
	}

	auth.pKey = config.Tokens.Key
	return *auth, nil
}

// Authorize implements gorilla/mux middleware that validates token in authorizatino header defined
// in config.Authorization. If token is invalid, function returns 404 for user.
// If token is valid, user_id is appended to request context
//
func (a *Auth) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := response.NewHttp(w)
		key := r.Header.Get(config.AuthorizationHeader)
		if key == "" {
			err := response.Unauthorized(resp)
			if err != nil {
				err = Error.Wrap(&err, "Failed to send response to client")
				logrus.Error(err)
			}
			return
		}

		parts := strings.Split(key, " ")
		if len(parts) != 2 {
			response.BadRequest("Invalid token", resp)
			return
		}

		next.ServeHTTP(w, r)
	})
}
