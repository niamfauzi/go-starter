package http

import (
	"fmt"
	"net/http"
	"time"
)

func NewServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
}
