package core

import (
	"fmt"
	"net/http"
	"strings"
)

const port = "8090"

//Start API server
func initAPI() {
	fmt.Println("Start API server on port " + port)
	go func() {
		http.HandleFunc("/", receivedURL)
		http.ListenAndServe(":"+port, nil)
	}()
}

//for HEALTHCHECK Dockerfile instruction
func receivedURL(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(404)
	list := strings.Split(req.Host, ".")
	if len(list) >= 2 {
		fmt.Fprintf(resp, "no server found for stack=%s service=%s\n", list[1], list[0])
		return
	}
	fmt.Fprintf(resp, "no server found for host: %s\n", req.Host)
}
