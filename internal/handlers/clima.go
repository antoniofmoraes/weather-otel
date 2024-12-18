package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/antoniofmoraes/weather/internal/services"
	"github.com/gorilla/mux"
)

func ClimaHandler(w http.ResponseWriter, r *http.Request) {
	cep := mux.Vars(r)["cep"]

	// Verifica se o CEP tem o formato correto
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	if !match {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Busca a localização a partir do CEP
	cidade, err := services.BuscarLocalizacaoPorCEP(cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	// Busca o clima atual da cidade
	clima, err := services.BuscarClima(cidade)
	if err != nil {
		http.Error(w, "error fetching weather data", http.StatusInternalServerError)
		return
	}

	response := map[string]float64{
		"temp_C": clima.Celsius,
		"temp_F": clima.Fahrenheit,
		"temp_K": clima.Kelvin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
