package auth

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"poi-service/cmd/download"
)

// JwkStore will take care to synchronize memory and remote backends to provide JWKs.
type JwkStore interface {
	// GetJWK will try to return the rawJWK from either local memory cache or from remote backends. The following errors
	// can occur:
	// - InvalidParameter: The given parameters are invalid. Not authorized.
	// - NoKeyAvailable: The rawJWK was retrieved successfully but contains not the requested JWK. Not authorized.
	// - other errors: Something goes wrong. Check the error/logs for details. Retry needed.
	GetJWK(kid, iss string) (rawJWK string, err error)
}

//------------------------------------------------------------------------------

// NoKeyAvailable indicates that there is no issue but also no JWK available at all.
const NoKeyAvailable = NoKeyAvailableError("no JWK available")

type NoKeyAvailableError string

func (e NoKeyAvailableError) Error() string { return string(e) }

//------------------------------------------------------------------------------

// InvalidParameter is given if request parameters do not fulfill required properties.
const InvalidParameter = InvalidParameterError("requested parameters invalid")

type InvalidParameterError string

func (e InvalidParameterError) Error() string { return string(e) }

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------

// DependencyMissing indicates that there is no issue but also no JWK available at all.
const DependencyMissing = DependencyMissingError("dependency missing")

type DependencyMissingError string

func (e DependencyMissingError) Error() string { return string(e) }

//------------------------------------------------------------------------------

// NewJwkStore creates a new cache instance.
// trustedBackendUrl: url that is trusted as iss
// client: http download client
func NewJwkStore(trustedBackendUrl string, client download.HttpRequester, cache JwkCache) JwkStore {
	store := jwkStore{trustedBackendUrl: trustedBackendUrl, client: client, cache: cache}
	return &store
}

// jwkStore implements interface JwkStore
type jwkStore struct {
	trustedBackendUrl string
	client            download.HttpRequester
	cache             JwkCache
}

func (j *jwkStore) GetJWK(kid, iss string) (rawJWKs string, err error) {
	// First try to get the jwk from cache -> we might downloaded it already
	rawJWKs, err = j.getFromCache(kid, iss)
	if err == nil {
		return
	}

	// not available we must download new jwks -> stores a found jwk to cache
	if err := j.fetchFromBackend(kid, iss); err != nil {
		return "", err
	}

	// Now get the downloaded jwk from cache
	rawJWKs, err = j.getFromCache(kid, iss)
	return
}

func (j *jwkStore) fetchFromBackend(kid, iss string) (err error) {
	log.Info().Msg("Downloading JWKs")

	if iss == "" || kid == "" {
		err = InvalidParameter
		log.Warn().Err(err).Str("kid", kid).Str("issuer", iss).Msg("GetJWK request failed")
		return
	}

	if j.client == nil {
		err = DependencyMissing
		log.Warn().Err(err).Msg("http client")
		return
	}

	// check if the backend is a trusted one -> otherwise someone can just its own server
	if iss != j.trustedBackendUrl {
		log.Warn().Err(err).Msg("iss not a trusted backend")
		return
	}

	rawData, err := j.client.GetContent(iss + ".well-known/jwks.json")
	if err != nil {
		return
	}

	var jwks Jwks
	err = json.Unmarshal([]byte(rawData), &jwks)
	if err != nil {
		log.Warn().Err(err).Msg("during unmarshal")
		return
	}

	// Store all found jwk in cache
	for _, val := range jwks.Keys {
		if len(val.Iss) == 0 {
			// add iss since it is needed as key in cache
			val.Iss = iss
		}
		j.cache.Add(val)
	}

	log.Info().Msg("Downloading JWKs done")
	return nil
}

func (j *jwkStore) getFromCache(kid, iss string) (rawJWKs string, err error) {
	jwk, err := j.cache.Get(kid, iss)
	if err != nil {
		return "", NoKeyAvailable
	}

	return jwk.String(), nil
}
