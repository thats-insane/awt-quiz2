package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hi im cahlil and i study information technology")
	})

	http.ListenAndServe(":3000", nil)
}
