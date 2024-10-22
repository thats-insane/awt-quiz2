package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("start", "method", r.Method, "path", r.URL.Path)
		defer slog.Info("end", "method", r.Method, "path", r.URL.Path)

		next(w, r)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hi im cahlil")
}

func aboutHander(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "i study information technology")
}

func main() {
	http.HandleFunc("/home", logRequest(homeHandler))
	http.HandleFunc("/about", logRequest(aboutHander))

	http.ListenAndServe(":3000", nil)
}