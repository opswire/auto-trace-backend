package chat

import "time"

type Chat struct {
	Id        int64
	BuyerId   int64
	SellerId  int64
	AdId      int64
	CreatedAt time.Time
}

type Message struct {
	Id        int64
	ChatId    int64
	SenderId  int64
	Text      string
	IsRead    bool
	CreatedAt time.Time
	Mine      bool
}
