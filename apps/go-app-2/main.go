package main

import (
	"context"
	"log"

	"github.com/marcelobiao/poc-opentelemetry/apps/go-app-2/app"
)

func main() {
	// config exporter
	exporter := app.GetZipkinExporter("http://zipkin:9411/api/v2/spans")
	// exporter := app.GetJaegerExporter("http://jaeger:14268/api/traces")

	// init otel
	tp := app.InitOtel(exporter)

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// run webservers
	app.StartFiberWebServer()
	// app.StartGinWebServer()
}
