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
	// initialize the dependencies
	topwordServiceClient := client.NewClient(getEnv("TOP_WORD_SERVICE", "http://localhost:8080"))
	topwordHandler := topword.NewHandler(topwordServiceClient)

	// sticking with the builtin http mux to have 0 dependencies
	mux := http.NewServeMux()
	mux.HandleFunc("/top-words", topwordHandler.FindTopWords)

	// get the address from env or use the default one;
	// can be replaced with viper or envconfig packages; but remember 0 dependencies
	address := getEnv("ADDRESS", "0.0.0.0:8081")
	log.Printf("listening@%s\n", address)

	// Start the server and listen for requests
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalln("failed to start server")
	}
}

// getEnv returns the value of a the given env. If it is empty or not defined then returns the default value specified
func getEnv(key string, def string) string {
	v := os.Getenv(strings.ToUpper(key))
	if v == "" {
		return def
	}
	return v
}
