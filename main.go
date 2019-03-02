package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var (
	port = getEnv("SERVER_PORT", "8080")
)

type request struct {
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	req := &request{
		Method:  r.Method,
		Headers: r.Header,
		URL:     r.URL.String(),
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Body = string(body)

	res, err := json.MarshalIndent(req, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	log.Println(res)
}

func main() {
	log.Printf("Listening on port: %s ...\n", port)
	addr := net.JoinHostPort("", port)
	http.HandleFunc("/", handleRequest)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
