package server

import (
	"fmt"

	"net/http"
)

func Run() error {
	port := 7540
	http.Handle("/", http.FileServer(http.Dir("web")))

	fmt.Println("Server is listening...")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
