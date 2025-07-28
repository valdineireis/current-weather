package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"` // Nome da cidade
	Erro       bool   `json:"erro"`       // Campo que indica se o CEP foi encontrado
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func getCityFromCEP(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return "", fmt.Errorf("falha na requisição para ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ViaCEP retornou status não esperado: %s", resp.Status)
	}

	var viaCEPResponse ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEPResponse); err != nil {
		return "", fmt.Errorf("falha ao decodificar a resposta do ViaCEP: %w", err)
	}

	// A API ViaCEP retorna "erro: true" se o CEP não for encontrado
	if viaCEPResponse.Erro {
		return "", fmt.Errorf("cep não encontrado na base do ViaCEP")
	}

	if viaCEPResponse.Localidade == "" {
		return "", fmt.Errorf("cidade não encontrada para o CEP fornecido")
	}

	return viaCEPResponse.Localidade, nil
}

func getWeatherForCity(city, apiKey string) (float64, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, encodedCity)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("falha na requisição para WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("WeatherAPI retornou status não esperado: %s", resp.Status)
	}

	var weatherResponse WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return 0, fmt.Errorf("falha ao decodificar a resposta da WeatherAPI: %w", err)
	}

	return weatherResponse.Current.TempC, nil
}
