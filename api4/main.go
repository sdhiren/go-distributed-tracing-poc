package main

import (
	"fmt"
	"log"
	"net/http"
	"tracing/tracelib"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func main() {
	if err := tracelib.InitializeTracing("Service4", "http://localhost:14268/api/traces"); err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}

	r := gin.Default()
	r.GET("/dong", func(c *gin.Context) {

		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
		ctx := otel.GetTextMapPropagator().Extract(c, propagation.HeaderCarrier(c.Request.Header))
		_, span := tracelib.StartSpanFromContext(ctx, spanName)

		defer span.End()

		c.JSON(http.StatusOK, gin.H{"message": "response from service4: dong"})
	})

	log.Fatal(r.Run(":8083"))
}
