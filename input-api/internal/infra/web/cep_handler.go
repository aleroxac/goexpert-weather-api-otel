package web

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aleroxac/goexpert-weather-api-otel/configs"
	"github.com/aleroxac/goexpert-weather-api-otel/input-api/internal/entity"
	"github.com/aleroxac/goexpert-weather-api-otel/input-api/internal/infra/repo"
	"github.com/aleroxac/goexpert-weather-api-otel/input-api/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type WebCEPHandler struct {
	CEPRepository entity.CEPRepositoryInterface
	Configs       *configs.Conf
	Tracer        trace.Tracer
}

func NewWebCEPHandler(conf *configs.Conf, tracer trace.Tracer) *WebCEPHandler {
	return &WebCEPHandler{
		CEPRepository: repo.NewCEPRepository(conf.OrchestratorApiHost, conf.OrchestratorApiPort),
		Configs:       conf,
		Tracer:        tracer,
	}
}

func (h *WebCEPHandler) Get(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	_, span := h.Tracer.Start(ctx, "INPUT_API:GET_CEP_TEMP")
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "fail to read the response", http.StatusInternalServerError)
		return
	}

	var cep_data usecase.CEPInputDTO
	err = json.Unmarshal(resp, &cep_data)
	if err != nil {
		http.Error(w, "fail to parse the cep_data", http.StatusInternalServerError)
		return
	}

	validate_cep_dto := usecase.ValidateCEPInputDTO{
		CEP: cep_data.CEP,
	}

	validateCEP := usecase.NewValidateCEPUseCase(h.CEPRepository)
	is_valid := validateCEP.Execute(validate_cep_dto)
	if !is_valid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	get_cep_dto := usecase.CEPInputDTO{
		CEP: cep_data.CEP,
	}

	getCEP := usecase.NewGetCEPUseCase(h.CEPRepository)
	err = getCEP.Execute(get_cep_dto)
	if err != nil {
		http.Error(w, "error getting cep", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	span.End()
}
