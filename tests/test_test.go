package tests

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/antoniofmoraes/weather/internal/handlers"
	"github.com/antoniofmoraes/weather/internal/services"
	"github.com/antoniofmoraes/weather/internal/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestClimaHandler_Success(t *testing.T) {
	LoadEnv()
	req := httptest.NewRequest(http.MethodGet, "/clima/01001000", nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://viacep.com.br/ws/01001000/json/",
		httpmock.NewStringResponder(200, `{"localidade": "Sao Paulo"}`))

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=Sao Paulo", os.Getenv("WEATHER_API_KEY")),
		httpmock.NewStringResponder(200, `{"current": {"temp_c": 25.0}}`))

	router := mux.NewRouter()
	router.HandleFunc("/clima/{cep}", handlers.ClimaHandler)
	router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Contains(t, rr.Body.String(), "temp_C")
	assert.Contains(t, rr.Body.String(), "temp_F")
	assert.Contains(t, rr.Body.String(), "temp_K")
}

func TestClimaHandler_InvalidZipcode(t *testing.T) {
	req, _ := http.NewRequest("GET", "/clima/abcd1234", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/clima/{cep}", handlers.ClimaHandler)
	router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)
	assert.Equal(t, "invalid zipcode\n", rr.Body.String())
}

func TestClimaHandler_ZipcodeNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/clima/99999999", nil)
	rr := httptest.NewRecorder()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://viacep.com.br/ws/99999999/json/",
		httpmock.NewStringResponder(404, ""))

	router := mux.NewRouter()
	router.HandleFunc("/clima/{cep}", handlers.ClimaHandler)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "can not find zipcode\n", rr.Body.String())
}

func TestConverterTemperaturas(t *testing.T) {
	temp := utils.ConverterTemperaturas(25.0)
	assert.Equal(t, 25.0, temp.Celsius)
	assert.Equal(t, 77.0, temp.Fahrenheit)
	assert.Equal(t, 298.0, temp.Kelvin)
}

func TestBuscarLocalizacaoPorCEP_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://viacep.com.br/ws/01001000/json/",
		httpmock.NewStringResponder(200, `{"localidade": "Sao Paulo"}`))

	cidade, err := services.BuscarLocalizacaoPorCEP("01001000")
	assert.NoError(t, err)
	assert.Equal(t, "Sao Paulo", cidade)
}

func TestBuscarLocalizacaoPorCEP_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://viacep.com.br/ws/99999999/json/",
		httpmock.NewStringResponder(404, ""))

	cidade, err := services.BuscarLocalizacaoPorCEP("99999999")
	assert.Error(t, err)
	assert.Equal(t, "", cidade)
}

func TestBuscarClima_Success(t *testing.T) {
	LoadEnv()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=Sao Paulo", os.Getenv("WEATHER_API_KEY")),
		httpmock.NewStringResponder(200, `{"current": {"temp_c": 25.0}}`))

	clima, err := services.BuscarClima("Sao Paulo")
	assert.NoError(t, err)
	assert.Equal(t, 25.0, clima.Celsius)
	assert.Equal(t, 77.0, clima.Fahrenheit)
	assert.Equal(t, 298.0, clima.Kelvin)
}

func TestBuscarClima_Error(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.weatherapi.com/v1/current.json?key=YOUR_WEATHER_API_KEY&q=Sao Paulo",
		httpmock.NewStringResponder(500, ""))

	clima, err := services.BuscarClima("Sao Paulo")
	assert.Error(t, err)
	assert.Equal(t, 0.0, clima.Celsius)
	assert.Equal(t, 0.0, clima.Fahrenheit)
	assert.Equal(t, 0.0, clima.Kelvin)
}

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}
}
