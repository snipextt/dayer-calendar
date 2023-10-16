package models

type Response struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
