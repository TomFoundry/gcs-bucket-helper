package server

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

// atheraLogin endpoint
func (s *Server) atheraLogin(w http.ResponseWriter, r *http.Request) {

	state := randToken()

	redirectURL := getLoginURL(state, s.atheraCfg, oauth2.SetAuthURLParam("audience", "https://public.elara.io"))

	fmt.Println("redirectURL: ", redirectURL)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (s *Server) atheraAuth(w http.ResponseWriter, r *http.Request) {

	if errVal, ok := r.URL.Query()["error"]; ok {
		errDescriptionVal := r.URL.Query()["error_description"]

		err := fmt.Errorf("%s: %s", errVal[0], errDescriptionVal[0])

		// InternalServerError because client errors (e.g. wrong password) are caught by Auth0 before ever reaching this point
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiKey := r.URL.Query()["code"][0]

	exchangeOpts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_verifier", s.atheraCfg.Verifier),
	}

	tok, err := s.atheraCfg.Conf.Exchange(r.Context(), apiKey, exchangeOpts...)
	if err != nil {
		log.Fatal("Failed exchange: ", err)
	}

	message := "You are now logged in to Athera. You can close this page and go back to the terminal.\n"
	formattedMessage := fmt.Sprintf("<html><body><p>%s</html></body></p>", message)
	w.Write([]byte(formattedMessage))

	go s.executor.ExecuteAthera(tok)
}
