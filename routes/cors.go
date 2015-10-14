package routes
import (
	"net/http"
	"log"
)

func CorsWrapper (h http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE")
		w.Header().Add("Access-Control-Allow-Headers","content-type")
		if (r.Method == "OPTIONS") {
			log.Print("options request")
			//handle preflight in here
		} else {
			h(w,r)
		}
	}
}
