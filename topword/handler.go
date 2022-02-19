package topword

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hysem/top-word-service/client"
	"github.com/hysem/top-word-service/topword"
)

// Handler implementation
type Handler struct {
	topWordServiceClient client.Client
}

// NewHandler returns an instance of Handler implementation
func NewHandler(topWordServiceClient client.Client) *Handler {
	return &Handler{
		topWordServiceClient: topWordServiceClient,
	}
}

// FindTopWords handles the request for finding the top words
func (h *Handler) FindTopWords(rw http.ResponseWriter, r *http.Request) {
	// Step 1: POST method allows sending far more data than the GET method so allowing only post method here
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request topword.FindTopWordsRequest
	// Step 2: Get the text from form field data
	request.Text = r.FormValue("text")

	// Step 3: Find the top words using the client for top word service
	topWords, err := h.topWordServiceClient.FindTopWords(r.Context(), &request)
	if err != nil {
		log.Println("failed to find top words", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Step 4: Marshal the results
	b, err := json.Marshal(topWords)
	// Step 5: Handle failure
	if err != nil {
		log.Println("failed to marshal top words", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Step 6: Write response
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}
