package chat

import (
	"context"
)

type Repository interface {
	StoreMessage(ctx context.Context, chatId int64, dto StoreMessageDTO) (Message, error)
	StoreChat(ctx context.Context, dto StoreChatDTO) (Chat, error)
	ListChats(ctx context.Context) ([]Chat, int64, error)
	ListMessagesByChatId(ctx context.Context, chatId int64) ([]Message, int64, error)
}

type Service struct {
	repository Repository
}

func NewService(
	repository Repository,
) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) StoreChat(ctx context.Context, dto StoreChatDTO) (Chat, error) {
	chat, err := s.repository.StoreChat(ctx, dto)
	if err != nil {
		return Chat{}, err
	}

	return chat, nil
}

func (s *Service) StoreMessage(ctx context.Context, chatId int64, dto StoreMessageDTO) (Message, error) {
	chat, err := s.repository.StoreMessage(ctx, chatId, dto)
	if err != nil {
		return Message{}, err
	}

	return chat, nil
}

func (s *Service) ListChats(ctx context.Context) ([]Chat, int64, error) {
	chats, count, err := s.repository.ListChats(ctx)
	if err != nil {
		return nil, 0, err
	}

	return chats, count, nil
}

func (s *Service) ListMessagesByChatId(ctx context.Context, chatId int64) ([]Message, int64, error) {
	chats, count, err := s.repository.ListMessagesByChatId(ctx, chatId)
	if err != nil {
		return nil, 0, err
	}

	return chats, count, nil
}
