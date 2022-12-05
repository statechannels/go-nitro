package apps

import "github.com/statechannels/go-nitro/types"

// NOTE: Would be nice to try to reuse objectives from protocols

type AppRequest struct {
	AppType     string            `json:"appType"`
	RequestType string            `json:"requestType"`
	ChannelId   types.Destination `json:"channelId"`

	Data interface{} `json:"data"`
}
