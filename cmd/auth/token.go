package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	les "github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
	"strconv"
)

type SignatureAlgorithm string

// Supported values for SignatureAlgorithm
const (
	RS256 SignatureAlgorithm = "RS256" // RSASSA-PKCS-v1.5 using SHA-256
)

const (
	// Kid key of the JWT header field that contains the key ID
	Kid = "kid"
	// Iss key of the JWT body field that contains the issuer
	Iss = "iss"
	// Expiration
	Exp = "exp"
)

// Token definition of operations executed on a JWT
type Token interface {
	// GetValueForHeaderKey returns the header value for the provided key, if the key is not available nil is returned
	GetValueForHeaderKey(key string) *string
	// GetValueForClaim get claim value for the provided key
	GetValueForClaim(key string) *string
	// GetClaims returns all claims
	GetClaims() Claims
	// IsValid check if token is valid against the provided rawJWK
	IsValid(rawJWK string) error
	// SetClaims adds claims if not already existing or overwrites existing ones with the provided values
	SetClaims(claimsToAdapt Claims)
	// Sign signs the token with the provided private RSA key and returns it as string. If something fails empty string is returned.
	// If kid is provided it will be set in header.
	Sign(privateKeyPem string, alg SignatureAlgorithm, kid string) string
}

// token implements interface Token
type token struct {
	jwt *jwt.Token
}

// Claims type that uses the map[string]interface{} for JSON decoding
// This is the default claims type if you don't supply one
type Claims map[string]interface{}

// NewUnverifiedToken creates a new unverified Token and will return error if the rawToken could not be parsed as JWT
// base64Jwt - the base64 encoded Token jwt as provided in the header authorization field (without Bearer)
func NewUnverifiedToken(base64Jwt string) (Token, error) {
	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	jwt, _, err := parser.ParseUnverified(base64Jwt, claims)

	if err != nil {
		return nil, err
	}

	return &token{jwt: jwt}, nil
}

func (j *token) IsValid(rawJWK string) error {
	ks, err := jwk.ParseString(rawJWK)
	if err != nil {
		return err
	}

	tok, err := les.ParseString(j.jwt.Raw, les.WithKeySet(ks))
	if err != nil {
		return err
	}

	if err := les.Validate(tok); err != nil {
		return err
	}

	return nil
}

func getStringFromMap(key string, lookup map[string]interface{}) *string {
	resultIface, ok := lookup[key]
	if !ok {
		log.Error().Msgf("Map does not contain key: %s", key)
		return nil
	}

	var result string
	switch v := resultIface.(type) {
	case int:
		res, _ := resultIface.(int)
		result = strconv.Itoa(res)
	case float64:
		res, _ := resultIface.(float64)
		result = fmt.Sprintf("%.0f", res)
	case string:
		result, _ = resultIface.(string)
	default:
		log.Error().Msgf("Could not cast: %v to string", v)
		return nil
	}

	return &result
}

// GetValueForHeaderKey see Token.GetValueForHeaderKey
func (j *token) GetValueForHeaderKey(key string) *string {
	return getStringFromMap(key, j.jwt.Header)
}

// GetValueForClaim see Token.GetValueForClaim
func (j *token) GetValueForClaim(key string) *string {
	return getStringFromMap(key, j.jwt.Claims.(jwt.MapClaims))
}

// GetClaims see Token.GetClaims
func (j *token) GetClaims() Claims {
	return Claims(j.jwt.Claims.(jwt.MapClaims))
}

func (j *token) Sign(privateKeyPem string, alg SignatureAlgorithm, kid string) (signed string) {
	if privateKeyPem == "" {
		log.Error().Msgf("privateKeyPem empty")
		return
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPem))

	if err != nil {
		log.Err(err).Msgf("failed parsing PEM")
		return
	}

	j.jwt.Method = jwt.GetSigningMethod(string(alg))
	j.jwt.Header["alg"] = string(alg)

	if kid != "" {
		j.jwt.Header["kid"] = kid
	}

	tokenString, err := j.jwt.SignedString(signKey)
	if err != nil {
		log.Err(err).Msgf("signing token failed")
		return
	}

	return tokenString
}

func (j *token) SetClaims(claimsToAdapt Claims) {
	for claim, value := range claimsToAdapt {
		j.jwt.Claims.(jwt.MapClaims)[claim] = value
	}
}
