package service

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"zipcode-temperature-system-service-b/internal/dto"
)

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func NewMockTemperatureService(apiKey string, client *mockClient) *TemperatureService {
	return &TemperatureService{
		apiKey: apiKey,
		client: client,
	}
}

func TestGetCity_Success(t *testing.T) {
	mockResp := dto.ViaCEPResponse{Localidade: "São Paulo"}
	mockRespBody, _ := json.Marshal(mockResp)

	mockClient := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(mockRespBody)),
			}, nil
		},
	}

	service := NewMockTemperatureService("test-api-key", mockClient)
	city, err := service.GetCity("01001000")

	assert.NoError(t, err)
	assert.Equal(t, "São Paulo", city)
}

func TestGetCity_InvalidZipcode(t *testing.T) {
	mockClient := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"erro": true}`))),
			}, nil
		},
	}

	service := NewMockTemperatureService("test-api-key", mockClient)
	_, err := service.GetCity("19818")

	assert.Error(t, err)
	assert.Equal(t, "unexpected status code: 400", err.Error())
}

func TestGetTemperature_Success(t *testing.T) {
	mockResp := dto.WeatherResponse{
		Current: dto.CurrentWeather{TempC: 28.5},
	}
	mockRespBody, _ := json.Marshal(mockResp)

	mockClient := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(mockRespBody)),
			}, nil
		},
	}

	service := NewMockTemperatureService("test-api-key", mockClient)
	temp, err := service.GetTemperature("São Paulo")

	assert.NoError(t, err)
	assert.Equal(t, 28.5, temp)
}

func TestGetTemperature_APIError(t *testing.T) {
	mockClient := &mockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error": "location not found"}`))),
			}, nil
		},
	}

	service := NewMockTemperatureService("test-api-key", mockClient)
	_, err := service.GetTemperature("UnknownCity")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error in response of API")
}
