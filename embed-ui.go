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

func HostNitroUI(port uint) {
	fmt.Println("Hosting UI on port", port)

	staticSite, err := fs.Sub(fs.FS(staticSiteRaw), "packages/nitro-gui/dist")
	if err != nil {
		log.Fatalf("Error parsing static site: %s", err)
	}
	http.Handle("/", http.FileServer(http.FS(staticSite)))
}
