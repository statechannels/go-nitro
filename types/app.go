package types

type AppRequest struct {
	AppId       string      `json:"appType"`
	RequestType string      `json:"requestType"`
	ChannelId   Destination `json:"channelId"`

	Data interface{} `json:"data"`
}
