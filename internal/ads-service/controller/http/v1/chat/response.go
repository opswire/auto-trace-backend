package chat

import (
	"car-sell-buy-system/internal/ads-service/domain/chat"
	"time"
)

type ChatResponse struct {
	Id         int64     `json:"id"`
	BuyerId    int64     `json:"buyer_id"`
	SellerId   int64     `json:"seller_id"`
	AdId       int64     `json:"ad_id"`
	CreatedAt  time.Time `json:"created_at"`
	AdTitle    string    `json:"ad_title"`
	BuyerName  string    `json:"buyer_name"`
	SellerName string    `json:"seller_name"`
	IsBuyer    bool      `json:"is_buyer"`
}

func newChatResponse(chat chat.Chat) ChatResponse {
	return ChatResponse{
		Id:         chat.Id,
		BuyerId:    chat.BuyerId,
		SellerId:   chat.SellerId,
		AdId:       chat.AdId,
		CreatedAt:  chat.CreatedAt,
		AdTitle:    chat.AdTitle,
		BuyerName:  chat.BuyerName,
		SellerName: chat.SellerName,
		IsBuyer:    chat.IsBuyer,
	}
}

type ListChatResponse struct {
	Chats []ChatResponse `json:"chats"`
}

func newListChatsResponse(chats []chat.Chat) ListChatResponse {
	responses := make([]ChatResponse, 0, len(chats))

	for _, chs := range chats {
		responses = append(responses, newChatResponse(chs))
	}

	return ListChatResponse{
		Chats: responses,
	}
}

type MessageResponse struct {
	Id        int64     `json:"id"`
	ChatId    int64     `json:"chat_id"`
	SenderId  int64     `json:"sender_id"`
	Text      string    `json:"text"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	Mine      bool      `json:"mine"`
	ImageUrl  string    `json:"image_url"`
}

func newMessageResponse(message chat.Message) MessageResponse {
	return MessageResponse{
		Id:        message.Id,
		ChatId:    message.ChatId,
		SenderId:  message.SenderId,
		Text:      message.Text,
		IsRead:    message.IsRead,
		CreatedAt: message.CreatedAt,
		Mine:      message.Mine,
		ImageUrl:  message.ImageUrl,
	}
}

type ListMessageResponse struct {
	Messages []MessageResponse `json:"messages"`
}

func newListMessageResponse(messages []chat.Message) ListMessageResponse {
	responses := make([]MessageResponse, 0, len(messages))

	for _, msg := range messages {
		responses = append(responses, newMessageResponse(msg))
	}

	return ListMessageResponse{
		Messages: responses,
	}
}
