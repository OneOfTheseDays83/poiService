package auth

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Authorizer interface {
	Authorize(next http.Handler) http.Handler
}

type authorizer struct {
	jwkStore JwkStore
}

func NewAuthorizer(jwkStore JwkStore) Authorizer {
	return &authorizer{jwkStore: jwkStore}
}

func (a *authorizer) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt := getBearerToken(r.Header)

		// check for provided jwt
		if jwt == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Missing auth token")
			return
		}

		// parse token and validate
		token, err := NewUnverifiedToken(jwt)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Warn().Err(err).Msg("parsing token failed")
			return
		}

		kid := token.GetValueForHeaderKey(Kid)
		iss := token.GetValueForClaim(Iss)

		if kid == nil || iss == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("iss or kid not available in jwt")
			return
		}

		if a.jwkStore == nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Msg("jwkStore is nil")
			return
		}

		// get JWK for JWT
		rawJWK, err := a.jwkStore.GetJWK(*kid, *iss)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error().Err(err).Msg("GetJwk failed")
			return
		}

		// validate token
		if err := token.IsValid(rawJWK); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error().Err(err).Msg("Token not valid")
			return
		}

		expired, err := isTokenExpired(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if expired {
			w.WriteHeader(http.StatusUnauthorized)
			log.Info().Msg("JWT expired")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getBearerToken(header http.Header) string {
	auth := header.Get("Authorization")
	if auth == "" {
		log.Warn().Msg("no Authorization header")
		return ""
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		log.Warn().Msg("Could not find bearer token in Authorization header")
		return ""
	}

	return token
}

func isTokenExpired(token Token) (expired bool, err error) {
	// check expiration
	exp := token.GetValueForClaim(Exp)
	if exp == nil {
		log.Warn().Msg("Failed to get expiration")
		err = errors.New("failed to get exp")
		return
	}

	i, err := strconv.ParseInt(*exp, 10, 64)
	if err != nil {
		log.Warn().Msg("Failed to format expiration")
		return
	}
	tm := time.Unix(i, 0)

	if time.Now().After(tm) {
		expired = true
		return
	}

	return false, nil
}
