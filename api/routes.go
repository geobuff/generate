package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) ping(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("PING SUCCESSFUL"))
}

func (s *Server) createTrivia(writer http.ResponseWriter, request *http.Request) {
	trivia, err := s.service.CreateTrivia()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(trivia)
}

func (s *Server) regenerateTrivia(writer http.ResponseWriter, request *http.Request) {
	date := mux.Vars(request)["date"]
	trivia, err := s.service.RegenerateTrivia(date)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(trivia)
}
