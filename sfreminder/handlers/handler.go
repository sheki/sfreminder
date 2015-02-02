package handlers

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/catan", catanHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")

}

func catanHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Catan Rules")
}
