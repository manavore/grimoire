package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/manavore/grimoire/internal/components"
)

type Handlers struct {
	s3client *s3.Client
}

func NewHandlers() (*Handlers, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(os.Getenv("S3_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("S3_ACCESS_KEY"),
			os.Getenv("S3_SECRET_KEY"),
			"",
		)),
	)

	if err != nil {
		log.Fatalf("unable to load SKD config: %v", err)
	}

	s3c := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
	})

	return &Handlers{s3client: s3c}, err
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	components.Home().Render(r.Context(), w)
}

func (h *Handlers) FormFile(w http.ResponseWriter, r *http.Request) {
}

func (h *Handlers) UploadFile(w http.ResponseWriter, r *http.Request) {

}
