//go:build headful
// +build headful

package main

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed packages/nitro-gui/dist/*
var staticSiteRaw embed.FS

func main() {
	var staticSite fs.FS

	if &staticSiteRaw != nil {
		var err error
		staticSite, err = fs.Sub(fs.FS(staticSiteRaw), "packages/nitro-gui/dist")
		if err != nil {
			log.Fatalf("Error parsing static site: %s", err)
		}
	}

	startWith(staticSite)
}
