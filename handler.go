package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type WeatherHandler struct {
	apiKey string
}

func NewWeatherHandler(apiKey string) *WeatherHandler {
	return &WeatherHandler{apiKey: apiKey}
}

type FinalResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// isValidCEP valida se a string do CEP tem 8 dígitos numéricos.
func isValidCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}

func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273
}

// Trata as requisições HTTP.
func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")

	if !isValidCEP(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity) // 422
		return
	}

	city, err := getCityFromCEP(cep)
	if err != nil {
		log.Printf("Erro ao buscar cidade para o CEP %s: %v", cep, err)
		http.Error(w, "can not find zipcode", http.StatusNotFound) // 404
		return
	}

	tempC, err := getWeatherForCity(city, h.apiKey)
	if err != nil {
		log.Printf("Erro ao buscar clima para a cidade %s: %v", city, err)
		// Usamos 500 aqui, pois é uma falha interna (nossa ou da API de clima)
		http.Error(w, "could not retrieve weather data", http.StatusInternalServerError)
		return
	}

	response := FinalResponse{
		TempC: tempC,
		TempF: celsiusToFahrenheit(tempC),
		TempK: celsiusToKelvin(tempC),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(response)
}
