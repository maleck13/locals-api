package routes
import (
	"net/http"
	"encoding/json"
)

type health struct  {
	HEALTH string `json:"health"`
}

func Ping(wr http.ResponseWriter, req *http.Request){
	wr.Header().Set("Content-Type", "application/json");
	h:= &health{HEALTH:"ok"};
	json.NewEncoder(wr).Encode(h);
}