package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/marcelobiao/poc-opentelemetry/apps/go-app/infra/opentelemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/trace"
)

var Exporter *zipkin.Exporter

func main() {
	zipkinConn, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatal(err)
	}
	Exporter = zipkinConn

	app := fiber.New(
		fiber.Config{
			Immutable: true,
			BodyLimit: 50 * 1024 * 1024,
		},
	)
	app.Use(MyMiddleware())
	app.Get("/", homeHandler)
	if err := app.Listen("localhost:8000"); err != nil {
		panic(err)
	}
}

func MyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := uuid.New()
		provider := opentelemetry.GetTracer(Exporter, id.String())
		tracer := provider.Tracer("go-app")
		c.Locals("provider", provider)
		c.Locals("tracer", tracer)
		otelMiddleware := otelfiber.Middleware(otelfiber.WithTracerProvider(provider))
		err := otelMiddleware(c)
		return err
	}
}

func homeHandler(c *fiber.Ctx) error {
	// provider := c.Locals("provider").(*trace.TracerProvider)
	tracer := c.Locals("tracer").(trace.Tracer)
	time.Sleep(time.Millisecond * 50)

	_, processFile := tracer.Start(c.Context(), "process-file")
	processFile.SetAttributes(attribute.String("info", "top"))
	time.Sleep(time.Millisecond * 50)
	c.SendStatus(http.StatusAccepted)
	c.JSON("home")
	processFile.End()
	return nil
}

func testHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("test"))
}
