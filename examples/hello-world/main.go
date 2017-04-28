package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joshbetz/config"
)

func main() {
	c := config.New("config.json")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var value string
		err := c.Get("value", &value)

		if err != nil {
			fmt.Fprintf(w, "Error: %s\n", err)
			return
		}

		fmt.Fprintf(w, "Value: %s\n", value)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
