package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	hostnameAPIServer = "api.discoverer.local"
	addrPortPat       = regexp.MustCompile(`\d+$`)
)

type APIServer struct {
	mapping Mapping
	server  *http.Server
}

func NewAPIServer(mapping Mapping) *APIServer {
	return &APIServer{
		mapping: mapping,
	}
}

func (s *APIServer) Serve(addr string) error {
	if len(addr) == 0 {
		addr = ":0"
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return err
	}

	addr = lis.Addr().String()
	port, err := strconv.Atoi(string(addrPortPat.FindSubmatch([]byte(addr))[0]))
	if err != nil {
		log.Println(err)
		return err
	}
	s.mapping.Update(uint32(port), int32(os.Getpid()), hostnameAPIServer)

	s.server = &http.Server{
		Handler: s.createHandler(),
	}

	log.Printf("start API server on %s\n", addr)
	err = s.server.Serve(lis)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *APIServer) createHandler() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Any("/ping", s.handlePing)
	return e
}

func (s *APIServer) handlePing(c echo.Context) error {
	c.String(http.StatusOK, "pong")
	return nil
}
