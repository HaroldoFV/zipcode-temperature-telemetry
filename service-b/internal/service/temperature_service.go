package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"zipcode-temperature-system-service-b/internal/dto"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TemperatureService struct {
	apiKey string
	client HttpClient
}

func NewTemperatureService(apiKey string) *TemperatureService {
	return &TemperatureService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (s *TemperatureService) GetCity(cep string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var viaCEPResponse dto.ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEPResponse); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if viaCEPResponse.Erro {
		return "", fmt.Errorf("CEP not found")
	}

	return viaCEPResponse.Localidade, nil
}

func (s *TemperatureService) GetTemperature(city string) (float64, error) {
	cityEncoded := url.QueryEscape(city)
	urlFull := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", s.apiKey, cityEncoded)

	req, err := http.NewRequest("GET", urlFull, nil)
	if err != nil {
		return 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error in response of API: status %d", resp.StatusCode)
	}

	var weatherResponse dto.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return 0, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return weatherResponse.Current.TempC, nil
}
