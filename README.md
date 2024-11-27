# Zipcode Temperature System API

## Descrição

A **Zipcode Temperature System API** fornece uma interface para consultar a temperatura atual de uma cidade brasileira
com base no CEP. A API usa a [API WeatherAPI](https://www.weatherapi.com/) para obter dados de temperatura
e [ViaCEP](https://viacep.com.br/) para buscar dados de localização.

## Funcionalidades

- Buscar cidade com base no CEP.
- Consultar a temperatura atual da cidade.

## Pré-requisitos

- **Docker**.
- Conta e chave de API na [WeatherAPI](https://www.weatherapi.com/).

## Configuração

1. Crie um arquivo `.env` na raiz do projeto e adicione sua chave de API para o WeatherAPI:

    ```dotenv
    WEATHER_API_KEY=your_weather_api_key
    ```

## Como Usar

1. **Construa e execute a aplicação usando Docker Compose**:

    ```bash
    docker-compose up --build
    ```

2. As aplicações: A e B estarão disponível respectivamente em:
   ```http://localhost:8081 ```

   ``` http://localhost:8080```


3. Você pode acessar a interface do Zipkin em http://localhost:9411 para visualizar os traces.