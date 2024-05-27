package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleroxac/goexpert-weather-api-otel/configs"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
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

func TestCEPHandler(t *testing.T) {
	router := chi.NewRouter()

	configs := configs.MockConfig()
	tracer := otel.Tracer("test")

	cep_address := "01001001"
	mockRepo := new(MockCEPRepository)
	mockRepo.On("Get", cep_address).Return(nil)
	mockRepo.On("IsValid", cep_address).Return(true)

	handler := &WebCEPHandler{
		CEPRepository: mockRepo,
		Configs:       configs,
		Tracer:        tracer,
	}

	router.Post("/cep", handler.Get)
	body := []byte(`{"cep": "01001001"}`)
	req, err := http.NewRequest("POST", "/cep", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
