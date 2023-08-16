package main

import (
	"fmt"
	"github.com/inikotoran/high-available-server/handler"
	"net/http"
)

func main() {
	h := handler.NewHandler()
	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			h.Put(w, r)
		case http.MethodGet:
			h.Get(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
