package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bhupeshbhatia/go-agg-monitoring/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"
)

func ErrorStackTrace(err error) string {
	return fmt.Sprintf("%+v\n", err)
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router = setAuthenticationRoute(router)
	return router
}

func setAuthenticationRoute(router *mux.Router) *mux.Router {
	router.HandleFunc("/create-data", service.LoadDataInMongo).Methods("GET", "OPTIONS")
	router.HandleFunc("/sen-perday", service.PerDaySenVal).Methods("POST", "OPTIONS")
	router.HandleFunc("/sen-table", service.SenTable).Methods("POST", "OPTIONS")
	return router
}

func main() {
	err := godotenv.Load()
	if err != nil {
		err = errors.Wrap(err,
			".env file not found, env-vars will be read as set in environment",
		)
		log.Println(err)
	}

	// service.TestIfDataGenerated()

	router := initRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":8081", n)
}
