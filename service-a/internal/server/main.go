package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"io"
	"log"
	"net/http"
	"regexp"
	_ "zipcode-temperature-system-service-a/docs"
	"zipcode-temperature-system-service-a/internal/dto"
)

func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

// ValidateZipCode é o handler que valida o CEP e encaminha a resposta
// @Summary Valida e encaminha um CEP
// @Description Valida se o CEP contém 8 dígitos e encaminha para o Serviço B se for válido.
// @Accept  json
// @Produce  json
// @Param   zipcode  body  dto.ZipcodeRequest  true  "CEP"
// @Success 200 {string} string "Zipcode is valid and forwarded to Service B"
// @Failure 422 {string} string "invalid zipcode"
// @Router /validate-zipcode [post]
func ValidateZipCode(w http.ResponseWriter, r *http.Request) {
	var req dto.ZipcodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tracer := otel.Tracer("zipcode-temperature-system-service-a")
	ctx, span := tracer.Start(r.Context(), "ValidateZipCode")
	defer span.End()

	if isValidCEP(req.CEP) {
		url := "http://service-b:8080/temperature/" + req.CEP
		request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			http.Error(w, "failed to contact service B", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				http.Error(w, "failed to read response from service B", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "service B returned an error", resp.StatusCode)
		}
	} else {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
	}
}

func initTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	exporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// @title Zipcode Validation API
// @version 1.0
// @description Esta é uma API para validar e encaminhar CEPs.
// @host localhost:8081
// @BasePath /
func main() {
	tp, err := initTracer("zipcode-temperature-system-service-a")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() { _ = tp.Shutdown(context.Background()) }()

	r := chi.NewRouter()

	r.Post("/validate-zipcode", ValidateZipCode)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8081", r))
}
