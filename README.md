# current-weather-api-go

<details>
<summary>Descrição</summary>

Objetivo: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

Requisitos:

- O sistema deve receber um CEP válido de 8 digitos
- O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin.
- O sistema deve responder adequadamente nos seguintes cenários:
  - Em caso de sucesso:
    - Código HTTP: 200
    - Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
  - Em caso de falha, caso o CEP não seja válido (com formato correto):
    - Código HTTP: 422
    - Mensagem: invalid zipcode
  - ​​Em caso de falha, caso o CEP não seja encontrado:
    - Código HTTP: 404
    - Mensagem: can not find zipcode
- Deverá ser realizado o deploy no Google Cloud Run.

Detalhes:

- Utilize a API viaCEP (ou similar) para encontrar a localização que deseja consultar a temperatura: https://viacep.com.br/
- Utilize a API WeatherAPI (ou similar) para consultar as temperaturas desejadas: https://www.weatherapi.com/
- Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C \* 1,8 + 32
- Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273
  - Sendo F = Fahrenheit
  - Sendo C = Celsius
  - Sendo K = Kelvin
- Testes automatizados demonstrando o funcionamento.
- Deploy realizado no Google Cloud Run (free tier).
</details>

### Estrutura

```
weather-api-go/
├── .env # Arquivo para a chave da API (não versionado)
├── Dockerfile # Define a imagem da aplicação
├── docker-compose.yml # Orquestra o container para teste local
├── go.mod # Gerencia as dependências do projeto
├── go.sum
├── handler.go # Contém a lógica principal do handler HTTP
├── handler_test.go # Contém os testes automatizados para o handler
├── main.go # Ponto de entrada da aplicação (configura o servidor)
└── services.go # Funções para se comunicar com as APIs externas (ViaCEP, WeatherAPI)
```

### Execute os testes:

```bash
go test -v
```

### Suba o container com Docker Compose:

```bash
docker-compose up --build
```

### Teste as rotas em outro terminal:

- Sucesso:

```bash
curl http://localhost:8080/weather/01001000
# Resposta esperada: {"temp_C":21.0,"temp_F":69.8,"temp_K":294.0} (valores podem variar)
```

- CEP inválido:

```bash
curl -i http://localhost:8080/weather/123
# Resposta esperada: HTTP/1.1 422 Unprocessable Entity ... invalid zipcode
```

- CEP não encontrado:

```bash
curl -i http://localhost:8080/weather/99999999
# Resposta esperada: HTTP/1.1 404 Not Found ... can not find zipcode
```

### Deploy

```bash
gcloud run deploy current-weather-api-go \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="WEATHER_API_KEY=sua_chave_secreta_aqui"
```
