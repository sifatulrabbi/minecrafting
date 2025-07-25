package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sifatulrabbi/minecrafting/iprelay/internal"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/register-ip", func(w http.ResponseWriter, r *http.Request) {
		ip := internal.GetClientIP(r)
		fmt.Println("Request IP:", ip)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"messages": "Success!"}`)
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Starting server at port :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
