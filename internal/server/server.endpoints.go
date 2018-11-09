package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/athera-io/gcs-bucket-helper/internal/gcp"
	"golang.org/x/oauth2"
)

// login endpoint
func (s *Server) login(w http.ResponseWriter, r *http.Request) {

	state := randToken()

	redirectURL := s.getLoginURL(state)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (s *Server) getLoginURL(state string) string {

	challenge := verifierToCode(s.verifier)

	authCodeOpts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("response_type", "code"),
	}

	// State can be some kind of random generated hash string.
	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
	return s.conf.AuthCodeURL(state, authCodeOpts...)
}

func verifierToCode(verifier string) string {

	h := sha256.New()
	h.Write([]byte(verifier))

	hashed := h.Sum(nil)

	encoded := base64.StdEncoding.EncodeToString(hashed)

	// Convert to Base64URL by replacing URL incompatible chars
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)

	return encoded
}

// auth endpoint
func (s *Server) auth(w http.ResponseWriter, r *http.Request) {

	apiKey := r.URL.Query()["code"][0]

	exchangeOpts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
		oauth2.SetAuthURLParam("code_verifier", s.verifier),
	}

	tok, err := s.conf.Exchange(r.Context(), apiKey, exchangeOpts...)
	if err != nil {
		log.Fatal("Failed exchange: ", err)
	}

	client := s.conf.Client(oauth2.NoContext, tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)

	var userData gcp.UserData

	if err := json.Unmarshal(data, &userData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := "You are now logged in. You can close this page and go back to the terminal.\n"
	formattedMessage := fmt.Sprintf("<html><body><p>%s</html></body></p>", message)
	w.Write([]byte(formattedMessage))

	go gcp.Do(tok, userData)
}
