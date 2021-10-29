package controllers

import (
	"encoding/json"
	"github.com/solarlabsteam/dvpn-openwrt/services/dvpnconf"
	"io/ioutil"
	"net/http"
)

func GetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(dvpnconf.Config.DVPN)
}

func PostConfig(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	config, err := dvpnconf.ValidateAndUnmarshal(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	resp, err := dvpnconf.Config.PostConfig(config)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(resp)
}
