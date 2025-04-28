package chat

import (
	"car-sell-buy-system/internal/ads-service/domain/chat"
)

type StoreChatRequest struct {
	SellerId int64 `json:"seller_id"`
	AdId     int64 `json:"ad_id"`
}

func (r StoreChatRequest) ToDTO() chat.StoreChatDTO {
	return chat.StoreChatDTO{
		SellerId: r.SellerId,
		AdId:     r.AdId,
	}
}

type StoreMessageRequest struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (r StoreMessageRequest) ToDTO() chat.StoreMessageDTO {
	return chat.StoreMessageDTO{
		Text: r.Text,
	}
}
