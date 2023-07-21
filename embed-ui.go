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

func hostNitroUI(guiPort uint, rpcPort uint) {
	staticSite, err := fs.Sub(fs.FS(staticSiteRaw), "packages/nitro-gui/dist")
	if err != nil {
		log.Fatalf("Error parsing static site: %s", err)
	}

	fs := http.FileServer(http.FS(staticSite))
	http.Handle("/", fs)

	http.HandleFunc("/rpc-port", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", rpcPort)
	})

	serverAddress := fmt.Sprintf(":%d", guiPort)

	url := fmt.Sprintf("http://localhost:%d/", guiPort)

	fmt.Printf("Hosting UI at %s\n", url)
	http.ListenAndServe(serverAddress, nil)
}
