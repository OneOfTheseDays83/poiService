package auth

import (
	"errors"
	"github.com/rs/zerolog/log"
	"time"
)

// JwkCache will take care to store, retrieve and flush downloaded jwk
type JwkCache struct {
	cache map[string]entry
}

type entry struct {
	jwk       Jwk
	timestamp time.Time
}

func (c *JwkCache) Init() {
	c.cache = make(map[string]entry)
}

func (c *JwkCache) Add(jwk Jwk) {
	if len(jwk.Kid) == 0 {
		return
	}

	c.cache[jwk.Kid+jwk.Iss] = entry{jwk, time.Now()}
	log.Info().Str("kid+iss", jwk.Kid+jwk.Iss).Msg("added to cache")
}

func (c *JwkCache) Get(kid, iss string) (Jwk, error) {
	if val, ok := c.cache[kid+iss]; ok {
		return val.jwk, nil
	}
	return Jwk{}, errors.New("not found")
}

func (c *JwkCache) Flush(age time.Duration) {
	// yes, for golang it is possible to delete while iterating
	for key, val := range c.cache {
		if time.Now().Sub(val.timestamp) > age {
			log.Info().Str("kid+iss", key).Msg("flushing due to age")
			delete(c.cache, key)
		}
	}
}
