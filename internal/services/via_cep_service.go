package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

func BuscarLocalizacaoPorCEP(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CEP não encontrado")
	}

	var resultado ViaCEPResponse
	err = json.NewDecoder(resp.Body).Decode(&resultado)
	if err != nil {
		return "", err
	}

	return resultado.Localidade, nil
}
