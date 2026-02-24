package server

import (
	"log"
	"net/http"
	"time"

	v1 "github.com/GoEnterpricePlatform/goEP-core/cmd/api/server/v1"
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(port string) *HttpServer {
	apiV1 := v1.New()

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      apiV1,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serv := &HttpServer{
		server: server,
	}

	return serv
}

func (serv *HttpServer) Start() {
	log.Printf("Http server running http://localhost%s", serv.server.Addr)
	log.Fatal(serv.server.ListenAndServe())
}
