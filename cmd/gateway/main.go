package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/antoniofmoraes/weather-otel/internal/services"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
)

type ZipcodeRequest struct {
	CEP string `json:"cep"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func main() {
	loadEnv()

	tracerProvider := services.InitTracer(os.Getenv("ZIPKKIN_URL"), "weather-gateway")
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleZipcodeRequest)
	log.Println("Gateway running on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func loadEnv() {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	godotenv.Load(filepath.Join(dir, ".env"))
}

func validateZipcode(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

func handleZipcodeRequest(w http.ResponseWriter, r *http.Request) {
	var request ZipcodeRequest

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid request body"})
		return
	}

	if !validateZipcode(request.CEP) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
		return
	}

	tracer := otel.Tracer("gateway-service")
	ctx, span := tracer.Start(r.Context(), "Call Weather Service")
	defer span.End()

	resp, err := callWeatherService(ctx, request.CEP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "error calling Weather Service"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	w.Write(responseBody)
	defer resp.Body.Close()
}

func callWeatherService(ctx context.Context, cep string) (*http.Response, error) {
	tracer := otel.Tracer("gateway-service")
	_, span := tracer.Start(ctx, "Weather Service Request")
	defer span.End()

	weatherServiceUrl := os.Getenv("WEATHER_API_URL")

	url := fmt.Sprintf("%s/clima/%s", weatherServiceUrl, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
