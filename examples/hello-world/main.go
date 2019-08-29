package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joshbetz/config"
)

func main() {
	c, err := config.New("config.json")
	if err != nil {
		log.Fatal(err)
	}

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
