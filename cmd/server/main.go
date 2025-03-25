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

	h, err := handlers.NewHandlers()

	if err != nil {
		logger.Error("unable to load handlers")
	}

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		h.Home(w, r)
	})

	router.HandleFunc("GET /file", func(w http.ResponseWriter, r *http.Request) {
		h.UploadFile(w, r)
	})

	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	logger.Info("Server started on :8080")
	srv.ListenAndServe()

}
