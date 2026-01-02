package main

import "github.com/reno1r/weiss/apps/service/internal/http"

func main() {
	server := http.NewServer()
	server.Listen(":8080")
}
