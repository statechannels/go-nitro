//go:build !embed_ui

package main

import "fmt"

func HostNitroUI(port uint) {
	fmt.Println("Not hosting UI.")
}
