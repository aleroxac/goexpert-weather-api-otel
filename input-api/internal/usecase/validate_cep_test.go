package usecase_test

import (
	"testing"

	"github.com/aleroxac/goexpert-weather-api-otel/configs"
	"github.com/aleroxac/goexpert-weather-api-otel/input-api/internal/infra/web"
	"github.com/aleroxac/goexpert-weather-api-otel/input-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

func TestValidateCEP(t *testing.T) {
	configs := configs.MockConfig()
	tracer := otel.Tracer("test")

	get_weather_dto := usecase.ValidateCEPInputDTO{
		CEP: "01001001",
	}
	webCEPHandler := web.NewWebCEPHandler(configs, tracer)

	validate_cep := usecase.NewValidateCEPUseCase(webCEPHandler.CEPRepository)
	weather_output := validate_cep.Execute(get_weather_dto)
	assert.True(t, weather_output)
}
