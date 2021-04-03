package main

import (
	"net/http"
	"strings"
)

func dashHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")

	if path != "" && path != "index" {
		http.NotFound(w, r)
		return
	}

	var c client
	if err := c.readCookie(w, r); err != nil {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return
	}

	err := c.updateStocks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.updateTransactions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := templates.ExecuteTemplate(w, "index.html", c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
