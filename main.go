package main
import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/maleck13/locals-api/routes"
	"database/sql"
	"github.com/maleck13/locals-api/data"
)


func main(){
	//set up db
	var db *sql.DB;
	db = data.DataBaseConnection();
	defer db.Close();


	//setup routes
	router :=SetUpRoutes();
	http.Handle("/", router)
	http.ListenAndServe(":9005",nil)
}


func SetUpRoutes()*mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/health",routes.CorsWrapper(routes.Ping)).Methods("GET","OPTIONS")
	router.HandleFunc("/profile",routes.CorsWrapper(routes.CreateProfile)).Methods("POST","OPTIONS")
	router.HandleFunc("/profiles/{county}",routes.CorsWrapper(routes.ListProfiles)).Methods("GET","OPTIONS")
	router.HandleFunc("/profile/{id}",routes.CorsWrapper(routes.GetProfile)).Methods("GET","OPTIONS")
	router.HandleFunc("/profile/{id}",routes.CorsWrapper(routes.UpdateProfile)).Methods("PUT","OPTIONS")
	router.HandleFunc("/profile/{id}",routes.CorsWrapper(routes.DeleteProfile)).Methods("DELETE","OPTIONS")
	return router
}