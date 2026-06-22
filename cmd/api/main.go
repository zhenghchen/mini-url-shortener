package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/zhenghchen/mini-url-shortener/pkg/database"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func generateCode() string {

	code := make([]byte, 6)

	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}

func main() {
	ctx := context.Background()
	
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	store, err := database.NewStore(ctx, endpoint)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handleShorten(w, r, store)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRedirect(w, r, store)
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)



}

func handleShorten(w http.ResponseWriter, r *http.Request, store *database.Store) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.URL == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	code := generateCode()

	if err := store.SaveURL(r.Context(), code, body.URL); err != nil {
		fmt.Println("save error:", err)
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"code": code})


}

func handleRedirect(w http.ResponseWriter, r *http.Request, store *database.Store) {
	code := r.URL.Path[1:]
	url, err := store.GetURL(r.Context(), code)
	if err != nil || url == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
