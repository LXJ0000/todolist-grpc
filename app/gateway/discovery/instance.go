package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

type Server struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Weight  int64  `json:"weight"`
}

func (s *Server) perfix() string {
	if s.Version == "" {
		s.Version = "latest"
	}
	return fmt.Sprintf("/%s/%s/", s.Name, s.Version)
}

func (s *Server) Path() string {
	return fmt.Sprintf("%s%s", s.perfix(), s.Addr)
}

func UnmarshalServer(bytes []byte) (Server, error) {
	var server Server
	err := json.Unmarshal(bytes, &server)
	return server, err
}

func SplitPath(path string) (Server, error) {
	var server Server
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return server, errors.New("invalid path")
	}
	server.Addr = parts[len(parts)-1] // TODO
	return server, nil
}

func (s *Server) exist(l []resolver.Address, addr resolver.Address) bool {
	for _, a := range l {
		if a.Addr == addr.Addr {
			return true
		}
	}
	return false
}

// Remove helper function
func Remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
