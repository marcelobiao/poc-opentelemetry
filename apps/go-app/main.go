package main

import (
	"context"
	"log"

	"github.com/marcelobiao/poc-opentelemetry/apps/go-app/app"
)

func main() {
	// config exporter
	exporter := app.GetZipkinExporter("http://zipkin:9411/api/v2/spans")

	// init otel
	tp := app.InitOtel(exporter)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// init database
	app.StartMongo("mongodb://test:test@mongodb:27017")

	// run webservers
	go app.StartGinWebServer()
	app.StartFiberWebServer()
}
