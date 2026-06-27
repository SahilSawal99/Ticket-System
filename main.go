package main
import (
	"fmt"
	"log"
	"net/http"

	"github.com/sahilsawal99/ticket-system/internal/auth"
	"github.com/sahilsawal99/ticket-system/internal/controller"
)

func main() {
	mux := http.NewServeMux()
	registerRoutes(mux)
	startServer(mux)
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/auth/login", controller.LoginHandler)
	mux.HandleFunc("/auth/register", controller.RegisterHandler)
	mux.HandleFunc("/health", controller.HealthHandler)

	securedTickets := auth.AuthMiddleware(controller.TicketsHandler)
	mux.HandleFunc("/tickets", securedTickets)
	mux.HandleFunc("/tickets/", securedTickets)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","endpoint":"/health"}`))
}

func startServer(handler http.Handler) {
	fmt.Println("Server started; available endpoints: /, /health, /auth/login, /auth/register, /tickets")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
