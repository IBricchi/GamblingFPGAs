package server

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *HttpServer) basicAuth(realm string, creds map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok {
				h.basicAuthFailed(w, realm, http.StatusText(401))
				return
			}

			credPassHash, credUserOk := creds[user]
			if !credUserOk {
				h.basicAuthFailed(w, realm, "Invalid username")
				return
			}

			byteCredPassHash := []byte(credPassHash)
			bytePass := []byte(pass)
			if err := bcrypt.CompareHashAndPassword(byteCredPassHash, bytePass); err != nil {
				h.basicAuthFailed(w, realm, "Wrong password")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *HttpServer) basicAuthFailed(w http.ResponseWriter, realm string, msg string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	http.Error(w, msg, http.StatusUnauthorized)

	h.logger.Info("Auth failed: " + msg)
}
