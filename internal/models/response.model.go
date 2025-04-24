package models

type Response struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
