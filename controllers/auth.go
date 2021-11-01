package controllers

import (
	"encoding/json"
	"github.com/solarlabsteam/dvpn-openwrt/services/auth"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	loginDetails, err := auth.ValidateAndUnmarshal(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err = auth.Store.Login(loginDetails.Username, loginDetails.Password); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	json.NewEncoder(w).Encode(auth.Store)
}
