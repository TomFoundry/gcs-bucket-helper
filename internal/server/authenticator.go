package server

import (
	"golang.org/x/oauth2"
)

type Config struct {
	AppID           string
	AppPseudoSecret string
	Conf            *oauth2.Config
	Verifier        string
}
