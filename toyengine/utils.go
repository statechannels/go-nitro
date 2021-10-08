package main

import "fmt"

// ChannelId computes a toy channel id for a virtual channel
func ChannelId(joiner uint, proposer uint) string {
	return fmt.Sprintf("%v-%v", joiner, proposer)
}
