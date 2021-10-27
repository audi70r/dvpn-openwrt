package controllers

import (
	"encoding/json"
	"github.com/audi70r/gochecknat"
	"net/http"
)

func GetNATInfo(w http.ResponseWriter, r *http.Request) {
	info, err := gochecknat.GetNATInfo()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(info)
}
