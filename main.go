package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/geobuff/generate/api"
	"github.com/geobuff/generate/storage"
	"github.com/geobuff/generate/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	listenAddr := flag.String("listenAddr", ":8081", "the server address")
	flag.Parse()

	store, err := storage.NewPostgresStore(os.Getenv("CONNECTION_STRING"))
	if err != nil {
		panic(err)
	}

	rateLimiterMax, _ := strconv.ParseFloat(os.Getenv("RATE_LIMITER_MAX"), 64)
	allowedOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
	allowedMethods := strings.Split(os.Getenv("CORS_METHODS"), ",")
	allowedHeaders := strings.Split(os.Getenv("CORS_HEADERS"), ",")

	service := utils.NewService(store)
	server := api.NewServer(*listenAddr, rateLimiterMax, allowedOrigins, allowedMethods, allowedHeaders, service)
	fmt.Println("server running on port:", *listenAddr)
	log.Fatal(server.Start())
}
