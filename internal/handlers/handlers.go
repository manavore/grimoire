package handlers

import "net/http"

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test"))
}
