package chat

import (
	"car-sell-buy-system/internal/chats-service/domain/chat"
	"car-sell-buy-system/pkg/storage/local"
	"mime/multipart"
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
	ChatId int64                 `form:"chat_id"`
	Text   string                `form:"text"`
	Image  *multipart.FileHeader `form:"image"`
}

func (r StoreMessageRequest) ToDTO() chat.StoreMessageDTO {
	var image *local.UploadedFile
	if r.Image != nil {
		img, err := local.ConvertUploadedFile(r.Image)
		if err == nil {
			image = img
		}
	}

	return chat.StoreMessageDTO{
		Text:  r.Text,
		Image: image,
	}
}
