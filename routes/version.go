package routes

import (
	"encoding/json"
	"net/http"
)

type version struct {
	VERSION string `json:"version"`
}

func Version(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	h := &version{VERSION: "1.1.1"}
	json.NewEncoder(wr).Encode(h)
}
