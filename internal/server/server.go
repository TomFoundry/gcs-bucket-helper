package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Server ...
type Server struct {
	port            string
	appID           string
	appPseudoSecret string
	conf            *oauth2.Config
	verifier        string
}

// New ...
func New(port, appID, appPseudoSecret string) *Server {
	s := &Server{
		port:            port,
		appID:           appID,
		appPseudoSecret: appPseudoSecret,
	}

	redirectURL := fmt.Sprintf("http://127.0.0.1:%s/auth", s.port)

	s.conf = &oauth2.Config{
		ClientID:     s.appID,
		ClientSecret: s.appPseudoSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	s.verifier = randVerifier()

	return s
}

func randVerifier() string {

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

// Serve ...
func (s *Server) Serve() error {
	http.HandleFunc("/login", s.login)
	http.HandleFunc("/auth", s.auth)

	fmt.Println("- This utility will create a service account and storage bucket in Google Cloud Platform.")
	fmt.Println("- We need to log in to Google.")
	fmt.Println("- Please navigate to:")
	fmt.Printf("http://localhost:%s/login\n", s.port)

	return http.ListenAndServe(":"+s.port, nil)
}
