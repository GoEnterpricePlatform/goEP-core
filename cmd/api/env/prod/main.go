package main


import (
	"cmp"
	"os"

	"github.com/amorindev/go-tmpl/cmd/api/server"
)

func main() {
	hsp := cmp.Or(os.Getenv("HTTP_SERVER_PORT"), "8000")

	httpSrv := server.NewHttpServer(hsp)
	httpSrv.Start()
}