package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	fonction "web/placeholder_web/function"
)

func main() {
	http.HandleFunc("/", placeholder)
	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func placeholder(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	list := strings.Split(strings.Replace(url, "/", "", 1), "/")
	buffer, err := fonction.Do(list)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
