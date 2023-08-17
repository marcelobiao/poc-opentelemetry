package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
		c.SendStatus(http.StatusOK)
		c.JSON("pong")
		return nil
	})
	if err := app.Listen(":8091"); err != nil {
		panic(err)
	}
}
