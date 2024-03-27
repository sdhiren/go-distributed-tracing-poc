package main

import (
	"log"
	"tracing/api1/router"
	// "router"
)

func main() {
	r := router.InitRoutes()

	log.Fatal(r.Run(":8080"))
}