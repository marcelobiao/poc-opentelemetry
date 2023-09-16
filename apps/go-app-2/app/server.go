package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

func StartGinWebServer() {
	r := gin.Default()

	r.Use(otelgin.Middleware(""))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
	_ = r.Run(":8090")
}

func StartFiberWebServer() {
	app := fiber.New()

	app.Use(otelfiber.Middleware())

	app.Get("/ping", func(c *fiber.Ctx) error {
		tp := otel.GetTracerProvider()
		_, span1 := tp.Tracer("").Start(c.UserContext(), "stage1")
		defer span1.End()
		time.Sleep(time.Microsecond * 150)
		_, span2 := tp.Tracer("").Start(c.UserContext(), "stage2")
		time.Sleep(time.Microsecond * 50)
		span2.End()
		time.Sleep(time.Microsecond * 50)
		c.SendStatus(http.StatusOK)
		c.JSON("pong")
		return nil
	})
	if err := app.Listen(":8090"); err != nil {
		panic(err)
	}
}
