package main

import (
	"database/sql"
	"github.com/maleck13/locals-api/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/maleck13/locals-api/data"
	"github.com/maleck13/locals-api/routes"
	"net/http"
)

func main() {
	//set up db
	var db *sql.DB
	db = data.DataBaseConnection()
	defer db.Close()

	//setup routes
	router := SetUpRoutes()
	http.Handle("/", router)
	http.ListenAndServe(":9005", nil)
}

func SetUpRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", routes.CorsWrapper(routes.Ping)).Methods("GET", "OPTIONS")
	router.HandleFunc("/version", routes.CorsWrapper(routes.Version)).Methods("GET", "OPTIONS")

	router.HandleFunc("/profile", routes.CorsWrapper(routes.CreateProfile)).Methods("POST", "OPTIONS")
	router.HandleFunc("/profiles/{county}", routes.CorsWrapper(routes.ListProfiles)).Methods("GET", "OPTIONS")
	router.HandleFunc("/profile/{id}", routes.CorsWrapper(routes.GetProfile)).Methods("GET", "OPTIONS")
	router.HandleFunc("/profile/{id}", routes.CorsWrapper(routes.UpdateProfile)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/profile/{id}", routes.CorsWrapper(routes.DeleteProfile)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/profile/{id}/profilepic", routes.CorsWrapper(routes.UploadProfilePic)).Methods("POST", "OPTIONS")

	projectRouter := router.PathPrefix("/project").Subrouter()
	projectRouter.HandleFunc("/project", routes.CorsWrapper(routes.CreateProject)).Methods("POST", "OPTIONS")
	projectRouter.HandleFunc("/{id}", routes.CorsWrapper(routes.GetProject)).Methods("GET", "OPTIONS")

	router.Handle("/projects", routes.CorsWrapper(routes.ListProjects)).Methods("GET", "OPTIONS")

	return router
}
