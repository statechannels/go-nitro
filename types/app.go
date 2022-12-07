package types

type AppRequest struct {
	From Address `json:"from"`

	AppId       string      `json:"appId"`
	RequestType string      `json:"requestType"`
	ChannelId   Destination `json:"channelId"`

	Data interface{} `json:"data"`
}
