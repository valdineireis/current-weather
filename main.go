package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Carrega variáveis de ambiente do arquivo .env para desenvolvimento local
	godotenv.Load()

	// A porta é fornecida pelo Cloud Run através da variável de ambiente PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão para ambiente local
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("Variável de ambiente WEATHER_API_KEY não definida.")
	}

	weatherHandler := NewWeatherHandler(apiKey)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{cep}", weatherHandler.ServeHTTP)

	log.Printf("Servidor iniciado na porta %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
