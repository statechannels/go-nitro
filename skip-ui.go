//go:build !embed_ui

package main

import "fmt"

func hostNitroUI(port uint) {
	fmt.Println("Not hosting UI.")
}
