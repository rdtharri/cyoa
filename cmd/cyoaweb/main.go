package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rdtharri/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "Port to start server on.")
	filename := flag.String("story", "gopher.json", "JSON File that contains a story")
	flag.Parse()
	fmt.Printf("Using the story: %s", *filename)

	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(file)
	if err != nil {
		panic(err)
	}

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
