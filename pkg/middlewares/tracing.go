package middlewares

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func initTracer() func() {
	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost("localhost:6831")))
	if err != nil {
		log.Fatalf("Failed to create Jaeger exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("advertisement-service"),
		)),
	)
	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown TracerProvider: %v", err)
		}
	}
}

func InitTracer(e *echo.Echo) func() {
	fn := initTracer()
	e.Use(otelecho.Middleware("advertisement-service"))
	return fn
}
