package auth

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"poi-service/cmd/download"
	"testing"
)

func Test_jwkStore_GetJWK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClient := download.NewMockHttpRequester(ctrl)

	kid := "unique"
	iss := "http://test.de"
	jwkToTest := Jwk{Iss: iss, Kid: kid}

	t.Run("found in cache", func(t *testing.T) {
		cache := JwkCache{}
		cache.Init()
		cache.Add(jwkToTest)

		storeToTest := NewJwkStore("", httpClient, cache)

		jwk, err := storeToTest.GetJWK(kid, iss)
		assert.Nil(t, err)
		assert.Equal(t, jwkToTest.String(), jwk)
	})

	t.Run("download: iss empty", func(t *testing.T) {
		cache := JwkCache{}
		cache.Init()

		storeToTest := NewJwkStore("", httpClient, cache)

		_, err := storeToTest.GetJWK(kid, "")
		assert.NotNil(t, err)
		assert.Equal(t, err, InvalidParameter)
	})

	t.Run("download: iss not trusted", func(t *testing.T) {
		cache := JwkCache{}
		cache.Init()

		storeToTest := NewJwkStore("abc", httpClient, cache)

		_, err := storeToTest.GetJWK(kid, "def")
		assert.NotNil(t, err)
	})

	t.Run("download: success", func(t *testing.T) {
		cache := JwkCache{}
		cache.Init()
		jwks := Jwks{}
		jwks.Keys = append(jwks.Keys, jwkToTest)

		httpClient.EXPECT().GetContent(iss+".well-known/jwks.json").Times(1).Return(jwks.String(), nil)

		storeToTest := NewJwkStore(iss, httpClient, cache)
		jwk, err := storeToTest.GetJWK(kid, iss)
		assert.Nil(t, err)
		assert.Equal(t, jwkToTest.String(), jwk)
	})
}
