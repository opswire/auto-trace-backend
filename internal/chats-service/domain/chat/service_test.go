package chat_test

import (
	"car-sell-buy-system/internal/chats-service/domain/chat"
	"context"
	"errors"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*gomock.Controller, *chat.Service, *MockRepository) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockRepository(ctrl)
	service := chat.NewService(mockRepo)
	return ctrl, service, mockRepo
}

func TestService_StoreChat_Success(t *testing.T) {
	ctrl, svc, mockRepo := setup(t)
	defer ctrl.Finish()

	dto := chat.StoreChatDTO{}
	expected := chat.Chat{Id: 1}

	mockRepo.EXPECT().StoreChat(gomock.Any(), dto).Return(expected, nil)

	result, err := svc.StoreChat(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestService_StoreChat_Error(t *testing.T) {
	ctrl, svc, mockRepo := setup(t)
	defer ctrl.Finish()

	dto := chat.StoreChatDTO{}
	mockRepo.EXPECT().StoreChat(gomock.Any(), dto).Return(chat.Chat{}, errors.New("store error"))

	_, err := svc.StoreChat(context.Background(), dto)
	assert.Error(t, err)
}

func TestService_StoreMessage_Success(t *testing.T) {
	ctrl, svc, mockRepo := setup(t)
	defer ctrl.Finish()

	dto := chat.StoreMessageDTO{}
	expected := chat.Message{Id: 1}

	mockRepo.EXPECT().StoreMessage(gomock.Any(), int64(123), dto).Return(expected, nil)

	result, err := svc.StoreMessage(context.Background(), 123, dto)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestService_StoreMessage_Error(t *testing.T) {
	ctrl, svc, mockRepo := setup(t)
	defer ctrl.Finish()

	dto := chat.StoreMessageDTO{}
	mockRepo.EXPECT().StoreMessage(gomock.Any(), int64(123), dto).Return(chat.Message{}, errors.New("msg error"))

	_, err := svc.StoreMessage(context.Background(), 123, dto)
	assert.Error(t, err)
}

func TestService_ListChats_Success(t *testing.T) {
	ctrl, svc, mockRepo := setup(t)
	defer ctrl.Finish()

	expectedChats := []chat.Chat{{Id: 1}}
	expectedCount := int64(1)

	mockRepo.EXPECT().ListChats(gomock.Any()).Return(expectedChats, expectedCount, nil)

	chats, count, err := svc.ListChats(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedChats, chats)
	assert.Equal(t, expectedCount, count)
}

func TestService_ListChats_Error(t *testing.T) {
	ctrl, svc, mockRepo := setup(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().ListChats(gomock.Any()).Return(nil, int64(0), errors.New("list error"))

	_, _, err := svc.ListChats(context.Background())
	assert.Error(t, err)
}

func TestService_ListMessagesByChatId_Success(t *testing.T) {
	ctrl, svc, mockRepo := setup(t)
	defer ctrl.Finish()

	expectedMsgs := []chat.Message{{Id: 1}}
	expectedCount := int64(1)

	mockRepo.EXPECT().ListMessagesByChatId(gomock.Any(), int64(42)).Return(expectedMsgs, expectedCount, nil)

	msgs, count, err := svc.ListMessagesByChatId(context.Background(), 42)
	assert.NoError(t, err)
	assert.Equal(t, expectedMsgs, msgs)
	assert.Equal(t, expectedCount, count)
}

func TestService_ListMessagesByChatId_Error(t *testing.T) {
	ctrl, svc, mockRepo := setup(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().ListMessagesByChatId(gomock.Any(), int64(42)).Return(nil, int64(0), errors.New("query failed"))

	_, _, err := svc.ListMessagesByChatId(context.Background(), 42)
	assert.Error(t, err)
}
