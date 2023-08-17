package app

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func StartGinWebServer() {
	r := gin.Default()

	r.Use(otelgin.Middleware(""))

	r.GET("/todo", func(c *gin.Context) {
		collection := DBClient.Database("todo").Collection("todos")
		cur, _ := collection.Find(c.Request.Context(), bson.D{})

		results := make([]interface{}, 0)
		cur.All(c, &results)

		client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
		req, _ := http.NewRequestWithContext(c.Request.Context(), "GET", "http://go-app-2:8090/ping", nil)
		res, _ := client.Do(req)
		body, _ := io.ReadAll(res.Body)
		results = append(results, string(body))

		c.JSON(http.StatusOK, results)
	})
	_ = r.Run(":8080")
}

func StartFiberWebServer() {
	app := fiber.New()

	app.Use(otelfiber.Middleware())

	app.Get("/todo", func(c *fiber.Ctx) error {
		collection := DBClient.Database("todo").Collection("todos")
		cur, _ := collection.Find(c.Context(), bson.D{})

		results := make([]interface{}, 0)
		cur.All(c.Context(), &results)

		resp, _ := otelhttp.Get(c.Context(), "http://go-app-2:8090/ping")
		results = append(results, resp)

		c.SendStatus(http.StatusOK)
		c.JSON(results)
		return nil
	})
	if err := app.Listen(":8081"); err != nil {
		panic(err)
	}
}
