package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gapat/goMicro/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Handle("/AllowedCountry", handlers.GetAllowedCountry()).Methods("GET")

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", os.Getenv("PORT")),
		Handler: router,
	}

	log.Printf("Country IP Server is listening on port 8080")

	server.ListenAndServe()
}
