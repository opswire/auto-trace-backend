package chat

import "car-sell-buy-system/pkg/storage/local"

type StoreChatDTO struct {
	SellerId int64
	AdId     int64
}

type StoreMessageDTO struct {
	Text            string
	Image           *local.UploadedFile
	CurrentImageUrl string
}
