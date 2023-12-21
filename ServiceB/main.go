package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// Start the tracer and defer the Stop method.
	tracer.Start()
	defer tracer.Stop()

	// Create a gin.Engine
	r := gin.New()

	// Use the tracer middleware with your desired service name.
	r.Use(gintrace.Middleware("serviceB"))

	// Set up some endpoints.
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world!")
	})

	r.GET("/service", func(c *gin.Context) {
		ctx := c.Request
		fmt.Println(ctx)
		sctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(ctx.Header))
		if err != nil {
			c.String(500, "Internal Server Error")
			log.Fatalf("There was an error %s", err)
		}
		span := tracer.StartSpan("test.span", tracer.ChildOf(sctx))
		defer span.Finish()

		c.JSON(http.StatusOK, gin.H{
			"message": "Headers received",
		})
	})

	// And start gathering request traces.
	r.Run(":8081")
}
