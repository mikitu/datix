package model

type Response struct {
	HttpResponse
	Request Request     `json:"request"`
	ReplyTo string
}

type HttpResponse struct {
	Data    interface{} `json:"data"`
	Err     interface{} `json:"err"`
	Status	int `json:"status"`
}
