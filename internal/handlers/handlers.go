package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/manavore/grimoire/internal/services/s3"

	"github.com/manavore/grimoire/internal/components"
	"github.com/manavore/grimoire/internal/components/fileUpload"
)

type Handlers struct {
	S3Uploader s3.S3Uploader
}

func NewHandlers(S3Uploader s3.S3Uploader) *Handlers {
	return &Handlers{
		S3Uploader: S3Uploader,
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	components.Home().Render(r.Context(), w)
}

func (h *Handlers) FormFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération du fichier: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	bucketName := os.Getenv("S3_BUCKET_NAME")

	fileName, err := h.S3Uploader.UploadFile(r.Context(), file, header, bucketName)
	if err != nil {
		http.Error(w, "Erreur lors de l'upload du fichier: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Fichier '%s' téléchargé avec succès", fileName)
}

func (h *Handlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	fileUpload.FileUpload().Render(r.Context(), w)
}
