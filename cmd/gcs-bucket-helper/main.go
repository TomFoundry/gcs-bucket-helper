package main

import (
	"fmt"
	"log"

	"github.com/athera-io/gcs-bucket-helper/internal/athera"
	executorImpl "github.com/athera-io/gcs-bucket-helper/internal/executor/impl"
	"github.com/athera-io/gcs-bucket-helper/internal/gcp"
	"github.com/athera-io/gcs-bucket-helper/internal/server"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	// All package variables are set at build time
	googleAppID           string
	googleAppPseudoSecret string
	atheraAppID           string
	voyagerURL            string
	atheraAuthURL         string
	atheraTokenURL        string
)

func main() {

	atheraExecutor, err := athera.New(voyagerURL)

	if err != nil {
		log.Fatal("Failed instantiating Athera executor: ", err)
	}

	gcpExecutor := gcp.New()

	e := executorImpl.NewExecutor(atheraExecutor, gcpExecutor)

	port := "8001"

	serv := server.New(
		e,
		port,
		makeAtheraCfg(port, atheraAppID, atheraAuthURL, atheraTokenURL),
		makeGoogleCfg(port, googleAppID, googleAppPseudoSecret),
	)

	if err := serv.Serve(); err != nil {
		log.Fatal("HTTP Server has failed: ", err)
	}
}

func makeAtheraCfg(port, appID, authURL, tokenURL string) *server.Config {
	conf := &oauth2.Config{
		ClientID:    appID,
		RedirectURL: fmt.Sprintf("http://127.0.0.1:%s/%s", port, server.EndpointAtheraAuth),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	return &server.Config{
		AppID:    appID,
		Conf:     conf,
		Verifier: server.RandVerifier(),
	}
}

func makeGoogleCfg(port, appID, appPseudoSecret string) *server.Config {
	conf := &oauth2.Config{
		ClientID:     appID,
		ClientSecret: appPseudoSecret,
		RedirectURL:  fmt.Sprintf("http://127.0.0.1:%s/%s", port, server.EndpointGoogleAuth),
		Scopes: []string{
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return &server.Config{
		AppID:           appID,
		AppPseudoSecret: appPseudoSecret,
		Conf:            conf,
		Verifier:        server.RandVerifier(),
	}
}
