package suscriber

type Suscriber struct {
	UserName string `json:"UserName"`
	ChatId   int64  `json:"ChatId"`
	Status   bool   `json:"Status"`
}
