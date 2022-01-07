package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	staticJwt = "eyJraWQiOiJ0ZXN0IiwiYWxnIjoiUlMyNTYiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJodHRwczovL215dGVzdC8iLCJzY3AiOltdLCJ2ZXIiOiIxIiwiaXNzIjoic29tZW9uZSIsInR5cCI6IkFUIiwidG50IjoibWUiLCJleHAiOjE2MTY2NzUzMDAsImlhdCI6MTYxNjY3MTcwMCwicnQtaWQiOiI3MDUyMzk3OS0zMTM2LTRiOTQtOTYwNi0xN2M2YWFiZTM2YTgiLCJqdGkiOiJkNzhmZTk2Ny0wMmQxLTQ2MjAtODFhOC03ZTIyNTQyYmEwMzYiLCJzdHlwIjoiVDEifQ.JiTzTKHHLtCEHD24nEpLlkZRfBIt0_Dj5nJ6ob9cRDtPoiYVDk4eU3KVSprh34biDpqy97u6ArAS7uxw9afYhnYy4ShBuw0xnTWjYh2GG-TqNn5cp0JlkinWhaf8dm9jaHA8waneHIn5tiec3E3gJek0QQcOh790S87mxXaIm5eh0IQfdpNFpmLKtg1JS5G9ufro270mZw9bh5EjGP_lRabfNZCkeo2nozMq7pj40fxACnNICCiini_YCbBYRRnikOSjYR8mDqYy-Ko9g1DZCR2p6IF_ddeUuhqJbyYwuTsVuLKxRPaMothg9JIsWR2dnOyoa1qOfGk_46MTAACxdw"
	rsaJwk    = "{\"kty\":\"RSA\",\"e\":\"AQAB\",\"use\":\"sig\",\"kid\":\"test\",\"alg\":\"RS256\",\"n\":\"4DHtdnJ7H9UfSZbEP5v1UBDnCBuP9vG2Tq7B6TYV_0gphVDPUADrflZcCrBeWmxwljrzW6SJuJfvf_FgPL-6jFw4Ib2Ao52AQxBEG0x0BPSFrC11PgekYE_gbBvJAYkWGtC2rBOVBJx8RQPIXza05PIjOW2Z99La5cxROuIydxO2gpezudnlXljw9s47rBHf45tVqs8kGPjN7N91o4QqcSdGewY3Ly--uXNPCd20A-DWzYS0ZCW7UBgI0LG5ImLcB_DPBB7x9WgRDnbGfzJwPVoyEyNbWkLeqjogp1J7spfTpwEFP1XN6dYb_1CwBylE55l_2qlCCoM4A-EouEBKeQ\"}"

	// both keys are just generated for unit test and nowhere else used
	privateKey = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAyHz3lu8OmeHnh67NE5NfK5v1FEFyKy5k/JJc/Tx4XRtVBhr4\nBZMDcbP5OYoY6Bury9KTyqHe36RCYlbSA/E5LyEXYdyS8nBNaRtpM2km3nT6zj11\nUGJtXrOCvDlbBoHf8feypqMWKiAP0YpY8yDS/4hCck2OD8weM6BuQcn2J2tjLEZk\n/69vMCoKGK0Q7bIxuWJpvC6J12LiE1YyjAbDcZtGKdXRG6KbULJXmyDg+dNCgBlw\nBBp9WiFPBpdFZzqYX5vAB3JpRnBib8rRi3Fm8Ldq74lsb6ueee82Jd232/QkPAq4\nEiy/5DAEi4eUlHmEhj9scjFGECd+/SGHddgHHQIDAQABAoIBACilpXDVYMl0EoPg\nvbU1ULs/sE1+A06b5l+KsQ2qf+ColPFa8GP47V7VFTdEN05/pbH6LHqNnOkMnWTg\no02nT2etttbhaG18tUUVCJwiun2pi9vae/ljKzdi/6N3oWvNUwD0riS4tdqui2Z5\nPRV11zF1h7sy3BV51bmz0gbGkoBlGhnUKcucdFL/mfXni+KTBsl3FpMK9+IJ3sSV\n9/R3Ri8o3ME3chSE9Z2lRhPa43fXB8WrAubEjxujXE7OAsab9J33i/wx9pLq+gKA\nZ83PUpBhXETFvEHzxwIUBCB/Tgb2Fy5X5hhZv0duu/SdIrb+a8/odnA7GUSfp+mT\n4xSY0wUCgYEA9ZBYMw50oJx2abNHHdHv7jPGcJOkVYTToXyqRHAw1H4LE8HTCs10\nNNfPmLS4HmdXdwRNTw2Qy92V8zW2eaH8YBHP+/hAJjDcSBgeBtQVVEAp2Gbgwzzo\nCmKqn8WHnZ7NxUIES2gwt2xUpt43EPUTlPruw+0EjOdZSxveZMqbzbsCgYEA0QI2\nr0JH7dG5MgtkfUnnp2ldx166MyWnL4JeiyogmX5/Efp+Xc1XpfH2HURf8uYhjlEb\nIa4fsF1/DFE6ePL863zwZKx0Rmt1wkOXPtpxPXy+4h9/7hVrc8ODlh2vAmaIdNEo\niazbusLX96vKaOqJZAGvmWHeBwXCpsOnHkRGRQcCgYByJLFKsjp1+aR1B3dUHiSX\npYtlAsvNUJuKoccHXtrjut7tRRgTGmMcuP/vLHm08DZQxTgmOdkHWi18SohSS4Bj\nK4Rwy/kNh4KtJEC4zdZIPjb1NwTc26/EPA6xi4C5PHrLaR9T6c9TQ1Cp6/rOsAx1\nIJrhiYem81anOgIK+b6oRwKBgBmVYcg2HsPXhgnAJz7GyxpM5XO//p7AHyTLmnMC\nZxciyr8SoGEu/2mKoouWkQAUd0sKVn3a6HoYF7MURkoDxD22/13zVhBAmxt6VosV\nBgN2v47COFCWQp7a8cJwQ7nRfyZ9a67ef87uhq0EVDRfcQ3SvwHRXvIRHHB5Rn2H\n8eoFAoGBAO+S+6CiQzpHdWQwoxQnlR44KowpHa9X2iQREDBzWbShO3ykahyyeJjO\n85cCnnbbhfNF0TV1d7e+6WqfDQ5c0Pk8lRbrHUqM+KaH1Kfl1NsaiH+HKShxIOGc\nijlONnSLXAct+fcHv3Y8EUGQ8RF7Wnuck0IMjOgKhsQXULWi4sOp\n-----END RSA PRIVATE KEY-----"
	publicKey  = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyHz3lu8OmeHnh67NE5Nf\nK5v1FEFyKy5k/JJc/Tx4XRtVBhr4BZMDcbP5OYoY6Bury9KTyqHe36RCYlbSA/E5\nLyEXYdyS8nBNaRtpM2km3nT6zj11UGJtXrOCvDlbBoHf8feypqMWKiAP0YpY8yDS\n/4hCck2OD8weM6BuQcn2J2tjLEZk/69vMCoKGK0Q7bIxuWJpvC6J12LiE1YyjAbD\ncZtGKdXRG6KbULJXmyDg+dNCgBlwBBp9WiFPBpdFZzqYX5vAB3JpRnBib8rRi3Fm\n8Ldq74lsb6ueee82Jd232/QkPAq4Eiy/5DAEi4eUlHmEhj9scjFGECd+/SGHddgH\nHQIDAQAB\n-----END PUBLIC KEY-----\n"
)

func generateValidTokenAndKey(t *testing.T) (string, string) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)
	token := jwt.New()
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	jwkKey, err := jwk.New(key)
	require.Nil(t, err)
	jwkKey.Set(jwk.KeyIDKey, "testKey")
	signed, err := jwt.Sign(token, jwa.RS256, jwkKey)
	pubKey, err := jwk.New(key.PublicKey)
	pubKey.Set(jwk.KeyIDKey, "testKey")
	require.Nil(t, err)
	publicKey, err := json.Marshal(pubKey)
	require.Nil(t, err)
	return string(signed), string(publicKey)
}

func generateValidTokenWithInvalidTimestamp(t *testing.T) (string, string) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)
	token := jwt.New()
	token.Set(jwt.ExpirationKey, time.Now().Unix())
	jwkKey, err := jwk.New(key)
	require.Nil(t, err)
	jwkKey.Set(jwk.KeyIDKey, "testKey")
	signed, err := jwt.Sign(token, jwa.RS256, jwkKey)
	pubKey, err := jwk.New(key.PublicKey)
	pubKey.Set(jwk.KeyIDKey, "testKey")
	require.Nil(t, err)
	publicKey, err := json.Marshal(pubKey)
	require.Nil(t, err)
	return string(signed), string(publicKey)
}

func TestNewUnverifiedToken(t *testing.T) {
	// Invalid Token, rsaJwk
	jwt, err := NewUnverifiedToken("")
	require.Nil(t, jwt)
	require.NotNil(t, err)
	log.Err(err).Msg("Expected JWT could not be created")

	jwt, err = NewUnverifiedToken(staticJwt)
	require.NotNil(t, jwt)
	require.Nil(t, err)

	require.NotNil(t, jwt.GetValueForHeaderKey("kid"))
	require.Equal(t, "test", *jwt.GetValueForHeaderKey("kid"))

	require.Nil(t, jwt.GetValueForHeaderKey("123"))

	require.NotNil(t, jwt.GetValueForClaim("iss"))
	require.Equal(t, "someone", *jwt.GetValueForClaim("iss"))

	require.Nil(t, jwt.GetValueForClaim("123"))
}

func TestGetStringFromMap(t *testing.T) {
	lookup := make(map[string]interface{}, 1)
	lookup["234"] = 3
	require.NotNil(t, getStringFromMap("234", lookup))
	require.Equal(t, "3", *getStringFromMap("234", lookup))
}

func TestIsValid(t *testing.T) {
	jwt, err := NewUnverifiedToken(staticJwt)
	require.NotNil(t, jwt)
	require.Nil(t, err)

	// Invalid JWK
	require.NotNil(t, jwt.IsValid(""))

	// JWK can not be used for verification
	err = jwt.IsValid(rsaJwk)
	require.NotNil(t, err)
	log.Err(err).Msg("Could not verify JWK")

	tok, jwk := generateValidTokenWithInvalidTimestamp(t)
	jwt, err = NewUnverifiedToken(tok)
	require.Nil(t, err)
	require.NotNil(t, tok)

	err = jwt.IsValid(jwk)
	require.NotNil(t, err)

	tok, jwk = generateValidTokenAndKey(t)
	jwt, err = NewUnverifiedToken(tok)
	require.Nil(t, err)
	require.NotNil(t, tok)
}

func Test_token_GetClaims(t *testing.T) {

	wantClaimKeys := []string{"aud", "scp", "ver", "iss", "typ", "tnt", "exp", "iat", "rt-id", "jti", "styp"}

	jwt, err := NewUnverifiedToken(staticJwt)
	require.NotNil(t, jwt)
	require.Nil(t, err)

	claims := jwt.GetClaims()

	require.Equal(t, len(wantClaimKeys), len(claims))

	for _, entry := range wantClaimKeys {
		_, found := claims[entry]
		require.True(t, found)
	}
}

func Test_token_Sign(t *testing.T) {
	jwt, err := NewUnverifiedToken(staticJwt)
	require.NotNil(t, jwt)
	require.Nil(t, err)

	claims := make(map[string]interface{})
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	jwt.SetClaims(claims)

	out := jwt.Sign(privateKey, RS256, "MBB")
	require.NotNil(t, out)

	finalToken, errParse := jwtGo.Parse(out, func(t *jwtGo.Token) (interface{}, error) {
		ret, errParse := jwtGo.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		return ret, errParse
	})

	require.Nil(t, errParse)
	require.True(t, finalToken.Valid)
}
