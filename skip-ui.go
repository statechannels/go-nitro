//go:build !embed_ui

package main

import "fmt"

func hostNitroUI(uint, uint) {
	fmt.Println("Not hosting UI.")
}
