package main

import (
	"log"
	"net/http"
	"os"

	"github.com/antoniofmoraes/weather-otel/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	r := mux.NewRouter()
	r.HandleFunc("/clima/{cep}", handlers.ClimaHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("WEATHER_API_KEY") == "" {
		log.Fatal("WEATHER_API_KEY env variable is needed")
	}

	log.Printf("Iniciando o servidor na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
