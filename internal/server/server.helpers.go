package server

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

func getLoginURL(state string, cfg *Config, opts ...oauth2.AuthCodeOption) string {

	challenge := verifierToCode(cfg.Verifier)

	authCodeOpts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("response_type", "code"),
	}

	authCodeOpts = append(authCodeOpts, opts...)

	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	return cfg.Conf.AuthCodeURL(state, authCodeOpts...)
}

func verifierToCode(verifier string) string {

	h := sha256.New()
	h.Write([]byte(verifier))

	hashed := h.Sum(nil)

	encoded := base64.StdEncoding.EncodeToString(hashed)

	// Convert to Base64URL by replacing URL incompatible chars.
	// (Using base64.URLEncoding instead StdEndcoding returned error 400 because it does not replace "=" with "")
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)

	return encoded
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// RandVerifier generates a PKCE code verifier
func RandVerifier() string {

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"-._~"

	b := make([]byte, 64)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
