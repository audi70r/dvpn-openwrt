package node

import (
	"encoding/json"
	"net/http"
)

func (n *Node) StartNode(w http.ResponseWriter, r *http.Request) {
	if err := n.StartNodeStreamOutputToSocket(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(n)
}

// GetNode will return basic node information
func (n *Node) GetNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(n)
}

// KillNode will kill the running node process
func (n *Node) KillNode(w http.ResponseWriter, r *http.Request) {
	if err := n.Kill(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(n)
}
