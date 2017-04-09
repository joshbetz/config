package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joshbetz/config"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var value string

		err := config.Get("config.json", "value", &value)

		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		fmt.Fprintf(w, "Value: %s", value)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
