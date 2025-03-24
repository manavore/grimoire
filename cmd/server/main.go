package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/manavore/grimoire/internal/handlers"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := http.NewServeMux()

	h := handlers.NewHandlers()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		h.Home(w, r)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	logger.Info("Server started on :8080")
	srv.ListenAndServe()

}
