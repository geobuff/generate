package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/didip/tollbooth"
	"github.com/geobuff/generate/utils"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	listenAddr     string
	rateLimiterMax float64
	allowedOrigins []string
	allowedMethods []string
	allowedHeaders []string
	service        utils.IService
}

func NewServer(listenAddr string, rateLimiterMax float64, allowedOrigins []string, allowedMethods []string, allowedHeaders []string, service utils.IService) *Server {
	return &Server{
		listenAddr,
		rateLimiterMax,
		allowedOrigins,
		allowedMethods,
		allowedHeaders,
		service,
	}
}

func (s *Server) Start() error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})

	if err != nil {
		fmt.Println(err)
	}

	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	router := mux.NewRouter()
	router.HandleFunc("/", s.ping)
	router.HandleFunc("/api/trivia", sentryHandler.HandleFunc(s.createTrivia)).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", sentryHandler.HandleFunc(s.regenerateTrivia)).Methods("PUT")

	limiter := tollbooth.LimitHandler(tollbooth.NewLimiter(s.rateLimiterMax, nil), s.handler(router))
	return http.ListenAndServe(s.listenAddr, limiter)
}

func (s *Server) handler(router http.Handler) http.Handler {
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: s.allowedOrigins,
		AllowedMethods: strings.Split(os.Getenv("CORS_METHODS"), ","),
		AllowedHeaders: strings.Split(os.Getenv("CORS_HEADERS"), ","),
	})

	return corsOptions.Handler(router)
}
