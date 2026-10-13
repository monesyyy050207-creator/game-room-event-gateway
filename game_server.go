package gamechat

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Server struct {
	client SessionVerifier
}

func NewServer(client SessionVerifier) *Server {
	return &Server{client: client}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "POST required"})
		return
	}
	var response any
	var err error
	switch r.URL.Path {
	case "/room-token":
		var input TokenRequest
		err = json.NewDecoder(r.Body).Decode(&input)
		if err == nil {
			response, err = IssueRoomToken(r.Context(), s.client, input)
		}
	case "/game-event":
		var input GameEvent
		err = json.NewDecoder(r.Body).Decode(&input)
		if err == nil {
			response, err = RouteGameEvent(r.Context(), s.client, input)
		}
	default:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "route not found"})
		return
	}
	if err != nil {
		status := http.StatusBadRequest
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode >= 400 && apiErr.StatusCode < 500 {
			status = apiErr.StatusCode
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
