package repo

import (
	"testing"

	"github.com/aleroxac/goexpert-weather-api-otel/configs"
	"github.com/aleroxac/goexpert-weather-api-otel/input-api/internal/entity"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestValidateCep(t *testing.T) {
	configs := configs.MockConfig()

	t.Run("valid_cep", func(t *testing.T) {
		cep := entity.NewCEP(
			"01001001",
			"Praça da Sé",
			"lado par",
			"Sé",
			"São Paulo",
			"SP",
			"3550308",
			"1004",
			"11",
			"7107",
		)
		repo := NewCEPRepository(configs.OrchestratorApiHost, configs.OrchestratorApiPort)
		validation := repo.IsValid(cep.CEP)
		assert.True(t, validation)
	})

	t.Run("invalid_cep", func(t *testing.T) {
		cep := entity.NewCEP(
			"01001-001",
			"Praça da Sé",
			"lado par",
			"Sé",
			"São Paulo",
			"SP",
			"3550308",
			"1004",
			"11",
			"7107",
		)
		repo := NewCEPRepository(configs.OrchestratorApiHost, configs.OrchestratorApiPort)
		validation := repo.IsValid(cep.CEP)
		assert.False(t, validation)
	})
}

func TestGetCEP(t *testing.T) {
	configs := configs.MockConfig()

	cep := entity.NewCEP(
		"01001-001",
		"Praça da Sé",
		"lado par",
		"Sé",
		"São Paulo",
		"SP",
		"3550308",
		"1004",
		"11",
		"7107",
	)
	repo := NewCEPRepository(configs.OrchestratorApiHost, configs.OrchestratorApiPort)

	// Inicia o httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Configura o mock da resposta HTTP
	httpmock.RegisterResponder(
		"GET",
		"http://"+configs.OrchestratorApiHost+":"+configs.OrchestratorApiPort+"/cep/"+cep.CEP,
		httpmock.NewStringResponder(200, `{"message": "Success"}`),
	)

	err := repo.Get(cep.CEP)
	assert.NoError(t, err)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["GET http://"+configs.OrchestratorApiHost+":"+configs.OrchestratorApiPort+"/cep/"+cep.CEP])
}
