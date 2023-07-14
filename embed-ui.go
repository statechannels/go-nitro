//go:build embed_ui

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
)

//go:embed packages/nitro-gui/dist/*
var staticSiteRaw embed.FS

func hostNitroUI(port uint) {
	staticSite, err := fs.Sub(fs.FS(staticSiteRaw), "packages/nitro-gui/dist")
	if err != nil {
		log.Fatalf("Error parsing static site: %s", err)
	}

	http.Handle("/", http.FileServer(http.FS(staticSite)))
	serverAddress := fmt.Sprintf(":%d", port)

	url := fmt.Sprintf("http://localhost:%d/", port)
	// If we're using the WS transport it will already be running a server so we don't need to start a new one
	if isListening(serverAddress) {
		fmt.Printf("Using transport http server to host UI at %s\n", url)
		return
	}

	fmt.Printf("Hosting UI at %s\n", url)
	http.ListenAndServe(serverAddress, nil)
}

// isListening returns true if the server is already listening on the given server path
func isListening(address string) bool {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return true // Address is already in use
	}
	defer listener.Close()

	return false
}
