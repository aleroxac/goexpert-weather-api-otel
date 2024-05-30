package web

import (
	"encoding/json"
	"fmt"
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

	_, span := h.Tracer.Start(ctx, "SPAN_READ_BODY")
	resp, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("fail to read the response: %v", err), http.StatusInternalServerError)
		return
	}
	span.End()

	_, span = h.Tracer.Start(ctx, "SPAN_DECODE_RESPONSE")
	var cep_data usecase.CEPInputDTO
	err = json.Unmarshal(resp, &cep_data)
	if err != nil {
		http.Error(w, fmt.Sprintf("fail to parse the cep_data: %v", err), http.StatusInternalServerError)
		return
	}
	span.End()

	_, span = h.Tracer.Start(ctx, "SPAN_VALIDATE_CEP")
	validate_cep_dto := usecase.ValidateCEPInputDTO{
		CEP: cep_data.CEP,
	}

	validateCEP := usecase.NewValidateCEPUseCase(h.CEPRepository)
	is_valid := validateCEP.Execute(validate_cep_dto)
	if !is_valid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	span.End()

	_, span = h.Tracer.Start(ctx, "SPAN_GET_CEP_RESPONSE_TIME")
	get_cep_dto := usecase.CEPInputDTO{
		CEP: cep_data.CEP,
	}

	getCEP := usecase.NewGetCEPUseCase(h.CEPRepository)
	err = getCEP.Execute(get_cep_dto)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting cep: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	span.End()
}
