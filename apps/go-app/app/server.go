package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func StartGinWebServer() {
	r := gin.Default()

	r.Use(otelgin.Middleware(""))

	r.GET("/todo", func(c *gin.Context) {
		collection := DBClient.Database("todo").Collection("todos")

		cur, findErr := collection.Find(c.Request.Context(), bson.D{})
		if findErr != nil {
			c.AbortWithError(500, findErr)
			return
		}
		results := make([]interface{}, 0)
		curErr := cur.All(c, &results)
		if curErr != nil {
			c.AbortWithError(500, curErr)
			return
		}
		c.JSON(http.StatusOK, results)
	})
	_ = r.Run(":8082")
}

func StartFiberWebServer() {
	app := fiber.New()
	app.Use(otelfiber.Middleware())
	app.Get("/todo", func(c *fiber.Ctx) error {
		collection := DBClient.Database("todo").Collection("todos")
		cur, _ := collection.Find(c.Context(), bson.D{})

		results := make([]interface{}, 0)
		cur.All(c.Context(), &results)

		c.SendStatus(http.StatusAccepted)
		c.JSON(results)
		return nil
	})
	if err := app.Listen("localhost:8082"); err != nil {
		panic(err)
	}
}
