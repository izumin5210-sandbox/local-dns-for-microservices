package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		switch req.Host {
		case "foo.izumin.local":
			req.URL.Host = "localhost:8001"
		case "bar.izumin.local":
			req.URL.Host = "localhost:8002"
		case "baz.izumin.local":
			req.URL.Host = "localhost:8003"
		}
	}
	rp := &httputil.ReverseProxy{Director: director}
	server := http.Server{
		Addr:    os.Getenv("APP_HOST"),
		Handler: rp,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
