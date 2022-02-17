package topword

import (
	"encoding/json"
	"net/http"

	"github.com/hysem/top-word-service/client"
	"github.com/hysem/top-word-service/topword"
)

type Handler struct {
	topWordServiceClient client.Client
}

func NewHandler(topWordServiceClient client.Client) *Handler {
	return &Handler{
		topWordServiceClient: topWordServiceClient,
	}
}

func (h *Handler) FindTopWords(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request topword.FindTopWordsRequest
	request.Text = r.FormValue("text")

	topWords, err := h.topWordServiceClient.FindTopWords(r.Context(), &request)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(topWords)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}
