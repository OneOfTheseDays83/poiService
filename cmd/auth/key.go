package auth

import "encoding/json"

// Jwks is a slice of JWK
type Jwks struct {
	Keys []Jwk `json:"keys"`
}

// JWK definition
type Jwk struct {
	Use string `json:"use,omitempty"`
	Kty string `json:"kty,omitempty"`
	Kid string `json:"kid,omitempty"`
	Iss string `json:"iss,omitempty"`
	Crv string `json:"crv,omitempty"`
	Alg string `json:"alg,omitempty"`
	K   string `json:"k,omitempty"`
	X   string `json:"x,omitempty"`
	Y   string `json:"y,omitempty"`
	N   string `json:"n,omitempty"`
	E   string `json:"e,omitempty"`
	// -- Following fields are only used for private keys --
	// RSA uses D, P and Q, while ECDSA uses only D. Fields Dp, Dq, and Qi are
	// completely optional. Therefore for RSA/ECDSA, D != nil is a contract that
	// we have a private key whereas D == nil means we have only a public key.
	D  string `json:"d,omitempty"`
	P  string `json:"p,omitempty"`
	Q  string `json:"q,omitempty"`
	Dp string `json:"dp,omitempty"`
	Dq string `json:"dq,omitempty"`
	Qi string `json:"qi,omitempty"`
	// Certificates
	X5c []string `json:"x5c,omitempty"`
}

func (j *Jwk) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(b)
}

func (j *Jwks) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(b)
}
