package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("")

func main() {
	// Configurando exporter
	zipkinConn, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatal(err)
	}

	// Configurando tracer
	tp := initTracer(zipkinConn)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Configurando app fiber
	app := fiber.New()
	app.Use(otelfiber.Middleware())
	app.Get("/", homeHandler)
	if err := app.Listen("localhost:8000"); err != nil {
		panic(err)
	}
}

func initTracer(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("go-app"),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func homeHandler(c *fiber.Ctx) error {
	time.Sleep(time.Millisecond * 50)
	_, span := tracer.Start(c.UserContext(), "HomeHandler", oteltrace.WithAttributes(attribute.String("id", "testeteste")))
	defer span.End()
	span.SetAttributes(attribute.String("info", "texto"))
	time.Sleep(time.Millisecond * 50)
	c.SendStatus(http.StatusAccepted)
	c.JSON("home")
	return nil
}
