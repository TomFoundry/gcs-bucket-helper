package server

import (
	"fmt"
	"net/http"

	"github.com/athera-io/gcs-bucket-helper/internal/executor"
)

const (
	EndpointGoogleLogin = "google-login"
	EndpointGoogleAuth  = "google-auth"
	EndpointAtheraLogin = "athera-login"
	EndpointAtheraAuth  = "athera-auth"
)

// Server ...
type Server struct {
	executor  executor.Executor
	port      string
	atheraCfg *Config
	googleCfg *Config
}

// New ...
func New(ex executor.Executor, port string, atheraCfg *Config, googleCfg *Config) *Server {
	s := &Server{
		executor:  ex,
		port:      port,
		atheraCfg: atheraCfg,
		googleCfg: googleCfg,
	}

	return s
}

// Serve ...
func (s *Server) Serve() error {
	http.HandleFunc("/"+EndpointGoogleLogin, s.googleLogin)
	http.HandleFunc("/"+EndpointGoogleAuth, s.googleAuth)

	http.HandleFunc("/"+EndpointAtheraLogin, s.atheraLogin)
	http.HandleFunc("/"+EndpointAtheraAuth, s.atheraAuth)

	fmt.Println("- This utility will create a service account and storage bucket in Google Cloud Platform (GCP), then connect it to one of your groups in Athera.")
	fmt.Println("- We need to log in to Google to get permission to create the bucket.")
	fmt.Println("- Please navigate to:")
	fmt.Printf("http://localhost:%s/%s\n", s.port, EndpointGoogleLogin)

	return http.ListenAndServe(":"+s.port, nil)
}
