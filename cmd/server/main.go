package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/manavore/grimoire/internal/config"
	"github.com/manavore/grimoire/internal/handlers"
	"github.com/manavore/grimoire/internal/services/s3"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := http.NewServeMux()

	s3Client, err := config.NewS3Client()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du client S3: %v", err)
	}

	s3Service := s3.NewService(s3Client)

	h := handlers.NewHandlers(s3Service)

	if err != nil {
		logger.Error("unable to load handlers")
	}

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		h.Home(w, r)
	})

	router.HandleFunc("GET /file", func(w http.ResponseWriter, r *http.Request) {
		h.UploadFile(w, r)
	})

	router.HandleFunc("POST /file", func(w http.ResponseWriter, r *http.Request) {
		h.FormFile(w, r)
	})

	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	logger.Info("Server started on :8080")
	srv.ListenAndServe()

}
