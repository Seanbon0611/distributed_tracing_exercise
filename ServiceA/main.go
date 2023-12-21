package main

import (
	"fmt"
	"io"
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
	r.Use(gintrace.Middleware("serviceA"))

	// Set up some endpoints.
	r.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()
		if span, ok := tracer.SpanFromContext(ctx); ok {
			// Set tag
			span.SetTag("test_tag", "sucessful")
		}
		c.String(200, "Welcome to Gopherdog! ʕ◔ϖ◔ʔ")
	})

	//Distributed Tracing
	r.GET("/test", func(c *gin.Context) {
		ctx := c.Request.Context()
		if span, ok := tracer.SpanFromContext(ctx); ok {
			fmt.Printf("The SpanID is %v", span.Context().SpanID())
			fmt.Printf(" The TraceID is %v", span.Context().TraceID())
			url := "http://serviceb:8081/service"
			//Makes request to serviceB endpoint
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalf(" error making http request: %s\n", err)
			}
			req = req.WithContext(ctx)
			// Inject the span Context in the Request headers
			err = tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(req.Header))
			if err != nil {
				c.String(500, "Internal Server Error")
				log.Fatalf("the error was: %s", err)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatalf(" error making http request: %s\n", err)
			}

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Fatalf(" error making http request: %s\n", err)
			}
			c.JSON(200, string(resBody))
		}
	})

	r.Run(":8080")
}
