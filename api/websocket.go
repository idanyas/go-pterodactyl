package api

type WebsocketDetails struct {
	Token     string `json:"token"`
	SocketURL string `json:"socket"`
}

type WebsocketResponse struct {
	Data WebsocketDetails `json:"data"`
	URL  string           `json:"socket"`
}
