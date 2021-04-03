package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	var err error
	key, err = generateKey()
	if err != nil {
		panic(fmt.Errorf("Key generation failed: %v", err))
	}

	println("running")
	assets := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))
	http.Handle("/assets/", assets)
	http.HandleFunc("/", dashHandler)
	http.HandleFunc("/ledger", ledgerHandler)
	http.HandleFunc("/auth", authHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(fmt.Errorf("Failed to start HTTP server: %v", err))
	}
}
