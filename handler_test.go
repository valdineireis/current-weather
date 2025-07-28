package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherHandler(t *testing.T) {

	// Handler de teste com uma chave de API de mentira
	handler := NewWeatherHandler("fake-api-key")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{cep}", handler.ServeHTTP)

	// Cenário 1: CEP Inválido (muito curto)
	t.Run("invalid zipcode format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/weather/12345", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("handler retornou código de status errado: obteve %v, esperava %v",
				status, http.StatusUnprocessableEntity)
		}

		expected := "invalid zipcode\n"
		if rr.Body.String() != expected {
			t.Errorf("handler retornou corpo inesperado: obteve %v, esperava %v",
				rr.Body.String(), expected)
		}
	})

	// Cenário 2: CEP Inválido (com letras)
	t.Run("invalid zipcode format with letters", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/weather/abcdefgh", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("handler retornou código de status errado: obteve %v, esperava %v",
				status, http.StatusUnprocessableEntity)
		}
	})

	// Cenário 3: CEP não encontrado (retorno 404 da nossa API)
	// Este teste fará uma chamada real para ViaCEP, que retornará erro para um CEP que não existe.
	t.Run("zipcode not found", func(t *testing.T) {
		// CEP com formato válido mas que não existe
		req := httptest.NewRequest("GET", "/weather/99999999", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler retornou código de status errado: obteve %v, esperava %v",
				status, http.StatusNotFound)
		}

		expected := "can not find zipcode\n"
		if rr.Body.String() != expected {
			t.Errorf("handler retornou corpo inesperado: obteve %v, esperava %v",
				rr.Body.String(), expected)
		}
	})

	// Cenário 4: Sucesso
	// Nota: Este é um teste de integração, pois depende de serviços externos (ViaCEP e WeatherAPI)
	t.Run("success case", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/weather/01001000", nil)
		rr := httptest.NewRecorder()

		realApiKey := "SUA_CHAVE_REAL_DA_WEATHERAPI" // Substitua pela sua chave para rodar este teste específico
		if realApiKey == "SUA_CHAVE_REAL_DA_WEATHERAPI" {
			t.Skip("Pulando teste de sucesso. Substitua a API Key para executá-lo.")
		}

		successHandler := NewWeatherHandler(realApiKey)
		successMux := http.NewServeMux()
		successMux.HandleFunc("GET /weather/{cep}", successHandler.ServeHTTP)

		successMux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Fatalf("handler retornou código de status errado: obteve %v, esperava %v. Erro: %s",
				status, http.StatusOK, rr.Body.String())
		}

		var response FinalResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Não foi possível decodificar a resposta JSON: %v", err)
		}

		if response.TempC == 0 || response.TempF == 0 || response.TempK == 0 {
			t.Errorf("Temperaturas não podem ser zero. Resposta: %+v", response)
		}
	})
}
