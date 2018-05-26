package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
)

type ProxyServer struct {
	mapping Mapping
	server  *http.Server
}

func NewProxyServer(mapping Mapping) *ProxyServer {
	return &ProxyServer{
		mapping: mapping,
	}
}

func (s *ProxyServer) Serve(addr string) error {
	s.server = &http.Server{
		Addr:    addr,
		Handler: &httputil.ReverseProxy{Director: s.handle},
	}
	log.Printf("start proxy server on %s\n", addr)
	err := s.server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (s *ProxyServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *ProxyServer) handle(req *http.Request) {
	host, ok := s.mapping.Lookup(req.Host)
	if ok {
		log.Printf("map %s -> %s\n", req.Host, host)
		req.URL.Host = host
	}
	req.URL.Scheme = "http"
}
