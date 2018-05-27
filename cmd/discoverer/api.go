package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
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
	e.PATCH("/containers/:containerID/mappings", s.handleAddMappings)

	return e
}

func (s *APIServer) handlePing(c echo.Context) error {
	c.String(http.StatusOK, "pong")
	return nil
}

func (s *APIServer) handleAddMappings(c echo.Context) error {
	var req struct {
		Mappings []struct {
			Port     uint32 `json:"port" validate:"required"`
			Pid      int32  `json:"pid" validate:"required"`
			Hostname string `json:"hostname" validate:"required"`
		} `json:"mappings" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	cID := c.Param("containerID")
	out, err := exec.Command("docker", "inspect", "--format={{json .NetworkSettings.Ports}}", cID).Output()
	if err != nil {
		return err
	}

	networkSettings := map[string][]struct{ HostPort string }{}
	err = json.Unmarshal(out, &networkSettings)
	if err != nil {
		return err
	}

	for _, m := range req.Mappings {
		if list, ok := networkSettings[fmt.Sprintf("%d/tcp", m.Port)]; ok {
			for _, ns := range list {
				hostPort, err := strconv.Atoi(ns.HostPort)
				if err != nil {
					return err
				}
				s.mapping.Update(uint32(hostPort), m.Pid, m.Hostname)
			}
		}
	}

	return nil
}
