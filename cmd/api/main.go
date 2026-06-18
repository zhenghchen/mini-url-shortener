package main

import (
	"fmt"
	"net/http"
	"math/rand"
	"encoding/json"
)

var urlStore = make(map[string]string)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func generateCode() string {
	code := make([]byte, 6)

	for i := range code {

		code[i] = charset[rand.Intn(len(charset))]

	}
	
	return string(code)
}


func main() {
	
	http.HandleFunc("/", handleRedirect)
	http.HandleFunc("/shorten", handleShorten) 

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
	
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		URL string `json:url`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.URL == "" {
	
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return

	}

	code := generateCode()
	urlStore[code] = body.URL
	

	// this adds the header to the response saying that its a json object
	w.Header().Set("Content-Type", "application/json")
	
	// new encoder streams json into the response, encode turns the map that we need
	// to return to the redirect into json. 
	json.NewEncoder(w).Encode(map[string]string{"code":code})


	
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// need to strip the first /
	code := r.URL.Path[1:]
	url, ok := urlStore[code]
	
	if !ok {

		http.Error(w, "not found", http.StatusNotFound)
		return

	}
	http.Redirect(w, r, url, http.StatusFound)

}
