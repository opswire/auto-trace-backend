package chat

import (
	"car-sell-buy-system/pkg/storage/local"
	"context"
)

type Storage interface {
	Save(file *local.UploadedFile) (string, error)
}

type Repository interface {
	StoreMessage(ctx context.Context, chatId int64, dto StoreMessageDTO) (Message, error)
	StoreChat(ctx context.Context, dto StoreChatDTO) (Chat, error)
	ListChats(ctx context.Context) ([]Chat, int64, error)
	ListMessagesByChatId(ctx context.Context, chatId int64) ([]Message, int64, error)
}

type Service struct {
	repository Repository
	storage    Storage
}

func NewService(
	repository Repository,
	storage Storage,
) *Service {
	return &Service{
		repository: repository,
		storage:    storage,
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
	if dto.Image != nil {
		pth, err := s.storage.Save(dto.Image)
		if err != nil {
			return Message{}, err
		}

		dto.CurrentImageUrl = pth
	}

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

	newChats := make([]Chat, 0, len(chats))
	for _, chat := range chats {
		if chat.BuyerId == ctx.Value("userId").(int64) {
			chat.IsBuyer = true
		}
		newChats = append(newChats, chat)
	}

	return newChats, count, nil
}

func (s *Service) ListMessagesByChatId(ctx context.Context, chatId int64) ([]Message, int64, error) {
	chats, count, err := s.repository.ListMessagesByChatId(ctx, chatId)
	if err != nil {
		return nil, 0, err
	}

	return chats, count, nil
}
