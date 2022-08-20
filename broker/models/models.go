package models

type Data struct {
	Round   int    `json:"round"`
	Sender  int    `json:"sender"`
	Id      int    `json:"id"`
	Message string `json:"message"`
}
