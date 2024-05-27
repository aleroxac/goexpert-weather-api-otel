package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleroxac/goexpert-weather-api-otel/configs"
	"github.com/aleroxac/goexpert-weather-api-otel/orchestrator-api/internal/entity"
	"github.com/aleroxac/goexpert-weather-api-otel/orchestrator-api/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
)

type MockCEPRepository struct {
	mock.Mock
}

func (m *MockCEPRepository) Get(cep string) ([]byte, error) {
	args := m.Called(cep)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCEPRepository) Convert(data []byte) (*entity.CEP, error) {
	args := m.Called(data)
	return args.Get(0).(*entity.CEP), args.Error(1)
}

func (m *MockCEPRepository) IsValid(cep string) bool {
	args := m.Called(cep)
	return args.Bool(0)
}

type MockWeatherRepository struct {
	mock.Mock
}

func (m *MockWeatherRepository) Get(localidade string, apiKey string) ([]byte, error) {
	args := m.Called(localidade, apiKey)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockWeatherRepository) ConvertToWeatherResponse(data []byte) (*entity.WeatherResponse, error) {
	args := m.Called(data)
	return args.Get(0).(*entity.WeatherResponse), args.Error(1)
}

func (m *MockWeatherRepository) ConvertToWeather(data *entity.WeatherResponse) (*entity.Weather, error) {
	args := m.Called(data)
	return args.Get(0).(*entity.Weather), args.Error(1)
}

func TestCEPHandler(t *testing.T) {
	router := chi.NewRouter()

	configs := configs.MockConfig()
	tracer := otel.Tracer("test")

	cep_address := "01001001"
	cep_data := []byte(`{
		"cep": "01001-001",
		"logradouro": "Praça da Sé",
		"complemento": "lado par",
		"bairro": "Sé",
		"localidade": "São Paulo",
		"uf": "SP",
		"ibge": "3550308",
		"gia": "1004",
		"ddd": "11",
		"siafi": "7107"
	}`)

	mockRepo := new(MockCEPRepository)
	mockRepo.On("Get", cep_address).Return(cep_data, nil)
	mockRepo.On("Convert", cep_data).Return(&entity.CEP{
		CEP:         "01001-001",
		Logradouro:  "Praça da Sé",
		Complemento: "lado par",
		Bairro:      "Sé",
		Localidade:  "São Paulo",
		UF:          "SP",
		IBGE:        "3550308",
		GIA:         "1004",
		DDD:         "11",
		SIAFI:       "7107",
	}, nil)
	mockRepo.On("IsValid", cep_address).Return(true)

	weather_data := []byte(`{
		"coord": {
			"lon": -46.6361,
			"lat": -23.5475
		},
		"weather": [
			{
				"id": 803,
				"main": "Clouds",
				"description": "broken clouds",
				"icon": "04d"
			}
		],
		"base": "stations",
		"main": {
			"temp": 21.1,
			"feels_like": 21.35,
			"temp_min": 19.75,
			"temp_max": 24.14,
			"pressure": 1024,
			"humidity": 80
		},
		"visibility": 10000,
		"wind": {
			"speed": 3.6,
			"deg": 140
		},
		"clouds": {
			"all": 75
		},
		"dt": 1716126286,
		"sys": {
			"type": 1,
			"id": 8394,
			"country": "BR",
			"sunrise": 1716111334,
			"sunset": 1716150647
		},
		"timezone": -10800,
		"id": 3448439,
		"name": "São Paulo",
		"cod": 200
	}`)

	mockWeatherRepo := new(MockWeatherRepository)
	mockWeatherRepo.On("Get", "São Paulo", configs.OpenWeathermapApiKey).Return(weather_data, nil)
	mockWeatherRepo.On("ConvertToWeatherResponse", weather_data).Return(&entity.WeatherResponse{
		Main: entity.WeatherDetails{
			Temp: 25.0,
		},
	}, nil)
	mockWeatherRepo.On("ConvertToWeather", &entity.WeatherResponse{Main: entity.WeatherDetails{Temp: 25.0}}).Return(
		&entity.Weather{
			Celcius:    25.0,
			Fahrenheit: 77.0,
			Kelvin:     298.15,
		}, nil)

	handler := &WebCEPHandler{
		CEPRepository:     mockRepo,
		WeatherRepository: mockWeatherRepo,
		Configs:           configs,
		Tracer:            tracer,
	}

	router.Get("/cep/{cep}", handler.Get)
	req, err := http.NewRequest("GET", "/cep/01001001", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 200)

	var weather *usecase.WeatherOutputDTO
	err = json.Unmarshal(rr.Body.Bytes(), &weather)
	assert.NoError(t, err)

	assert.IsType(t, weather, &usecase.WeatherOutputDTO{})
	assert.NotEmpty(t, weather.Celcius)
	assert.NotEmpty(t, weather.Fahrenheit)
	assert.NotEmpty(t, weather.Kelvin)
}
