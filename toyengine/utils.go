package main

import "fmt"

func ChannelId(joiner uint, proposer uint) string {
	return fmt.Sprintf("%v-%v", joiner, proposer)
}
