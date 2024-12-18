package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/antoniofmoraes/weather-otel/internal/utils"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func BuscarClima(cidade string) (utils.Temperatura, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, cidade)

	resp, err := http.Get(url)
	if err != nil {
		return utils.Temperatura{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return utils.Temperatura{}, fmt.Errorf("Erro ao buscar o clima")
	}

	var resultado WeatherAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&resultado)
	if err != nil {
		return utils.Temperatura{}, err
	}

	tempC := resultado.Current.TempC
	return utils.ConverterTemperaturas(tempC), nil
}
