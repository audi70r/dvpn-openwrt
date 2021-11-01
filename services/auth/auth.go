package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/solarlabsteam/dvpn-openwrt/utilities/shadow"
	"net/http"
)

type LoginRequest struct {
	Username string
	Password string
}

type Auth struct {
	Token uuid.UUID
}

var Store Auth

func (s *Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authToken, authTokenDecodeErr := base64.StdEncoding.DecodeString(authHeader)
		if authTokenDecodeErr != nil {
			http.Error(w, authTokenDecodeErr.Error(), http.StatusUnauthorized)
			w.Write([]byte{})
			return
		}

		authTokenUUID, uuidParseErr := uuid.Parse(string(authToken))
		if uuidParseErr != nil {
			http.Error(w, uuidParseErr.Error(), http.StatusUnauthorized)
			w.Write([]byte{})
			return
		}

		if authTokenUUID != s.Token {
			http.Error(w, "user not authorized", http.StatusUnauthorized)
			w.Write([]byte{})
			return
		}

		next.ServeHTTP(w, r)
	})

	return nil
}

func (s *Auth) Login(username, password string) error {
	user, err := shadow.Lookup(username)
	if err != nil {
		return err
	}

	if err = user.VerifyPassword(password); err != nil {
		return err
	}

	if !user.IsPasswordValid() {
		return fmt.Errorf("password invalid")
	}

	s.Token = uuid.New()
	return nil
}

func (s *Auth) BasicAuthForHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			user, err := shadow.Lookup(username)

			if err != nil {
				w.Header().Add("Clear-Site-Data", "*")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if err = user.VerifyPassword(password); err != nil {
				w.Header().Add("Clear-Site-Data", "*")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if user.IsPasswordValid() {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func (s *Auth) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			user, err := shadow.Lookup(username)

			if err != nil {
				w.Header().Add("Clear-Site-Data", "*")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if err = user.VerifyPassword(password); err != nil {
				w.Header().Add("Clear-Site-Data", "*")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if user.IsPasswordValid() {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
