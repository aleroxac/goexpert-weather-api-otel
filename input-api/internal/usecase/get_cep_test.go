package usecase_test

import (
	"fmt"
	"testing"

	"github.com/aleroxac/goexpert-weather-api-otel/input-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCEPRepository struct {
	mock.Mock
}

func (m *MockCEPRepository) Get(cep string) error {
	args := m.Called(cep)
	return args.Error(0)
}

func (m *MockCEPRepository) IsValid(cep string) bool {
	args := m.Called(cep)
	return args.Bool(0)
}

func TestGetCEPUseCase(t *testing.T) {
	t.Run("valid cep", func(t *testing.T) {
		cep_address := "01001001"
		mockRepo := new(MockCEPRepository)
		mockRepo.On("Get", cep_address).Return(nil)

		get_cep_dto := usecase.CEPInputDTO{
			CEP: cep_address,
		}
		getCEP := usecase.NewGetCEPUseCase(mockRepo)
		err := getCEP.Execute(get_cep_dto)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid cep", func(t *testing.T) {
		cep_address := "010010"
		mockRepo := new(MockCEPRepository)
		mockRepo.On("Get", cep_address).Return(fmt.Errorf("invalid cep"))

		get_cep_dto := usecase.CEPInputDTO{
			CEP: cep_address,
		}
		getCEP := usecase.NewGetCEPUseCase(mockRepo)
		err := getCEP.Execute(get_cep_dto)
		assert.EqualError(t, err, "invalid cep")

		mockRepo.AssertExpectations(t)
	})
}
