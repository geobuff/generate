package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/didip/tollbooth"
	"github.com/geobuff/generate/utils"
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
	router := mux.NewRouter()
	router.HandleFunc("/", s.ping)
	router.HandleFunc("/api/trivia", s.createTrivia).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", s.regenerateTrivia).Methods("PUT")

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

func (s *Server) ping(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("PING SUCCESSFUL"))
}

func (s *Server) createTrivia(writer http.ResponseWriter, request *http.Request) {
	err := s.service.CreateTrivia()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) regenerateTrivia(writer http.ResponseWriter, request *http.Request) {
	date := mux.Vars(request)["date"]
	err := s.service.RegenerateTrivia(date)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
