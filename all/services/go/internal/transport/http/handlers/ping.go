package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type pingResponse struct {
	OK      bool   `json:"ok"`
	Service string `json:"service"`
	TS      int64  `json:"ts"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(pingResponse{
		OK:      true,
		Service: "api",
		TS:      time.Now().Unix(),
	})
}
