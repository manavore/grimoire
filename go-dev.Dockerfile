FROM golang:latest
RUN go install github.com/air-verse/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN git config --global --add safe.directory /app
WORKDIR /app
CMD ["air", "-c", ".air.toml"]