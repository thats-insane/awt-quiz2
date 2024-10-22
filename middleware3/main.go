package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("Request logged", "remoteAddr", r.RemoteAddr, "method", r.Method, "url", r.URL.RequestURI(), "duration", time.Since(start))
	})
}

func recoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

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

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "i study information technology")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/home", logRequest(homeHandler))
	mux.HandleFunc("/about", logRequest(aboutHandler))

	handler := securityHeadersMiddleware(recoverPanicMiddleware(logRequestMiddleware(mux)))

	http.ListenAndServe(":3000", handler)
}
