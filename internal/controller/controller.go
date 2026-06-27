package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sahilsawal99/ticket-system/internal/auth"
	"github.com/sahilsawal99/ticket-system/internal/service"
)

func TicketsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/tickets")
	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")

	switch {
	case path == "" || path == ".":
		if r.Method == http.MethodPost {
			createTicketHandler(w, r)
			return
		}
		if r.Method == http.MethodGet {
			listTicketsHandler(w, r)
			return
		}
	case len(parts) == 1 && r.Method == http.MethodGet:
		getTicketHandler(w, r, parts[0])
		return
	case len(parts) == 2 && parts[1] == "status" && r.Method == http.MethodPatch:
		updateTicketStatusHandler(w, r, parts[0])
		return
	}

	http.NotFound(w, r)
}

func createTicketHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ticket := service.CreateTicket(r.Header.Get("X-User-ID"), req.Title, req.Description)
	writeJSON(w, http.StatusCreated, ticket)
}

func listTicketsHandler(w http.ResponseWriter, r *http.Request) {
	tickets := service.ListTickets(r.Header.Get("X-User-ID"))
	writeJSON(w, http.StatusOK, tickets)
}

func getTicketHandler(w http.ResponseWriter, r *http.Request, id string) {
	ticket, err := service.GetTicket(r.Header.Get("X-User-ID"), id)
	if err != nil {
		switch err {
		case service.ErrTicketNotFound:
			http.Error(w, "Not found", http.StatusNotFound)
		case service.ErrForbidden:
			http.Error(w, "Forbidden", http.StatusForbidden)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, ticket)
}

func updateTicketStatusHandler(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ticket, err := service.UpdateTicketStatus(r.Header.Get("X-User-ID"), id, req.Status)
	if err != nil {
		switch err {
		case service.ErrTicketNotFound:
			http.Error(w, "Not found", http.StatusNotFound)
		case service.ErrForbidden:
			http.Error(w, "Forbidden", http.StatusForbidden)
		case service.ErrInvalidStatusTransition:
			http.Error(w, "Invalid status transition", http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := service.RegisterUser(req.Username, req.Password); err != nil {
		if err == service.ErrUsernameTaken {
			http.Error(w, "Username taken", http.StatusConflict)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	token, err := service.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return auth.AuthMiddleware(next)
}
