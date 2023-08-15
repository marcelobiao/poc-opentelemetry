package main

import (
	"context"
	"log"

	"github.com/marcelobiao/poc-opentelemetry/apps/go-app/app"
)

func main() {
	// Configurando exporter
	exporter := app.GetZipkinExporter("http://localhost:9411/api/v2/spans")

	// Configurando tracer
	tp := app.InitOtel(exporter)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Configurando database
	app.StartMongo()

	// Configurando server
	app.StartGinWebServer()
	// app.StartFiberWebServer()
}
