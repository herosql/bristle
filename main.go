package main

import (
	"net/http"

	"github.com/herosql/bristle/framework"
)

func main() {
	server := &http.Server{
		Handler: framework.NewCore(),
		Addr: ":8080",
	}
	server.ListenAndServe()
}
