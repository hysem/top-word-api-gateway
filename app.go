package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hysem/top-word-api-gateway/topword"
	"github.com/hysem/top-word-service/client"
)

const N = 10

func main() {
	topwordServiceClient := client.NewClient(getEnv("TOP_WORD_SERVICE", "http://localhost:8080"))
	topwordHandler := topword.NewHandler(topwordServiceClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/top-words", topwordHandler.FindTopWords)

	address := getEnv("ADDRESS", "0.0.0.0:8081")
	log.Printf("listening@%s\n", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalln("failed to start server")
	}
}

func getEnv(key string, def string) string {
	v := os.Getenv(strings.ToUpper(key))
	if v == "" {
		return def
	}
	return v
}
