package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var ErrUnauthorized = errors.New("unauthorized")

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

type Verifier struct {
	jwtSecret string
	ecKeys    map[string]*ecdsa.PublicKey
}

func NewVerifier(jwtSecret, supabaseURL string) *Verifier {
	v := &Verifier{
		jwtSecret: jwtSecret,
		ecKeys:    make(map[string]*ecdsa.PublicKey),
	}
	_ = v.fetchJWKS(supabaseURL)
	return v
}

func (v *Verifier) fetchJWKS(supabaseURL string) error {
	resp, err := http.Get(supabaseURL + "/auth/v1/.well-known/jwks.json") //nolint:gosec
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ks jwks
	if err := json.NewDecoder(resp.Body).Decode(&ks); err != nil {
		return err
	}

	for _, k := range ks.Keys {
		if k.Kty != "EC" || k.Crv != "P-256" || k.Kid == "" {
			continue
		}
		xBytes, err := base64.RawURLEncoding.DecodeString(k.X)
		if err != nil {
			continue
		}
		yBytes, err := base64.RawURLEncoding.DecodeString(k.Y)
		if err != nil {
			continue
		}
		v.ecKeys[k.Kid] = &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(xBytes),
			Y:     new(big.Int).SetBytes(yBytes),
		}
	}
	return nil
}

func (v *Verifier) Verify(tokenString string) error {
	_, err := jwt.Parse(tokenString, v.keyFunc)
	if err != nil {
		return ErrUnauthorized
	}
	return nil
}

func (v *Verifier) GetUserID(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, v.keyFunc)
	if err != nil {
		return "", ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnauthorized
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", ErrUnauthorized
	}

	return sub, nil
}

func (v *Verifier) keyFunc(token *jwt.Token) (any, error) {
	switch token.Method.(type) {
	case *jwt.SigningMethodHMAC:
		return []byte(v.jwtSecret), nil
	case *jwt.SigningMethodECDSA:
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, ErrUnauthorized
		}
		key, ok := v.ecKeys[kid]
		if !ok {
			return nil, ErrUnauthorized
		}
		return key, nil
	default:
		return nil, ErrUnauthorized
	}
}

func ExtractBearerToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return "", ErrUnauthorized
	}
	return strings.TrimPrefix(header, "Bearer "), nil
}
