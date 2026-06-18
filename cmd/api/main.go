package main

import (
	"fmt"
	"net/http"
)

func main() {
	
	http.HandleFunc("/", handleRedirect)
	http.HandleFunc("/shorten", handleShorten) 

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
	
}

func handleShorten(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "shorten endpoint - coming soon")

}

func handleRedirect(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "redirect endpoint - coming soon")

}
