package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

// googleLogin endpoint
func (s *Server) googleLogin(w http.ResponseWriter, r *http.Request) {

	state := randToken()

	redirectURL := getLoginURL(state, s.googleCfg)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// googleAuth endpoint
func (s *Server) googleAuth(w http.ResponseWriter, r *http.Request) {

	apiKey := r.URL.Query()["code"][0]

	exchangeOpts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_verifier", s.googleCfg.Verifier),
	}

	tok, err := s.googleCfg.Conf.Exchange(r.Context(), apiKey, exchangeOpts...)
	if err != nil {
		log.Fatal("Failed exchange: ", err)
	}

	client := s.googleCfg.Conf.Client(oauth2.NoContext, tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)

	var ud userData

	if err := json.Unmarshal(data, &ud); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := "You are now logged in to Google. You can close this page and go back to the terminal.\n"
	formattedMessage := fmt.Sprintf("<html><body><p>%s</html></body></p>", message)
	w.Write([]byte(formattedMessage))

	go s.executor.ExecuteGCP(tok, ud.Email, fmt.Sprintf("http://localhost:%s/%s\n", s.port, EndpointAtheraLogin))
}

type userData struct {
	Email string `json:"email"`
}
