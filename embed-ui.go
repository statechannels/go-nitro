//go:build embed_ui

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
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

	fmt.Printf("Hosting UI at %s\n", url)
	http.ListenAndServe(serverAddress, nil)
}
