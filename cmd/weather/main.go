package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/antoniofmoraes/weather-otel/internal/handlers"
	"github.com/antoniofmoraes/weather-otel/internal/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	tracerProvider := services.InitTracer(os.Getenv("ZIPKIN_URL"), "weather-gateway")
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("WEATHER_API_KEY") == "" {
		log.Fatal("WEATHER_API_KEY env variable is needed")
	}

	r := mux.NewRouter()
	r.HandleFunc("/clima/{cep}", handlers.ClimaHandler).Methods("GET")

	log.Printf("Iniciando o servidor na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
