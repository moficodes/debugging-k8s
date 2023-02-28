package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func root(w http.ResponseWriter, r *http.Request) {
	for _, e := range os.Environ() {
		fmt.Fprintln(w, e)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func badHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "NOT OK")
}

func oomKill(w http.ResponseWriter, r *http.Request) {
	var s []string
	for {
		s = append(s, "Hello, World!")
	}
}

func logger(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		f(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", logger(root))
	mux.HandleFunc("/healthz", logger(healthz))
	mux.HandleFunc("/badhealthz", logger(badHealthz))
	mux.HandleFunc("/oomkill", logger(oomKill))

	PORT := ":8080"
	log.Println("Server is running on port", PORT)

	log.Fatal(http.ListenAndServe(PORT, mux))
}
