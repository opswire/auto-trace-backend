package ad_test

import (
	"car-sell-buy-system/internal/ads-service/domain/ad"
	"car-sell-buy-system/internal/ads-service/domain/nft"
	"car-sell-buy-system/pkg/storage/local"
	"context"
	"errors"
	"go.uber.org/mock/gomock"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*gomock.Controller, *ad.Service, *MockRepository, *MockNftRepository, *MockStorage) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockRepository(ctrl)
	mockNftRepo := NewMockNftRepository(ctrl)
	mockStorage := NewMockStorage(ctrl)
	service := ad.NewService(mockRepo, mockNftRepo, mockStorage)
	return ctrl, service, mockRepo, mockNftRepo, mockStorage
}

func TestService_GetById_Success(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	expectedAd := ad.Ad{}
	mockRepo.EXPECT().GetById(gomock.Any(), int64(1)).Return(expectedAd, nil)

	result, err := svc.GetById(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedAd, result)
}

func TestService_GetById_Error(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().GetById(gomock.Any(), int64(1)).Return(ad.Ad{}, errors.New("not found"))

	_, err := svc.GetById(context.Background(), 1)
	assert.Error(t, err)
}

func TestService_Store_WithImage_Success(t *testing.T) {
	ctrl, svc, mockRepo, _, mockStorage := setup(t)
	defer ctrl.Finish()

	image := &local.UploadedFile{Name: "image.jpg"}
	dto := ad.StoreDTO{Image: image}

	mockStorage.EXPECT().Save(image).Return("path/to/image.jpg", nil)
	mockRepo.EXPECT().Store(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input ad.StoreDTO) (ad.Ad, error) {
		assert.Equal(t, "path/to/image.jpg", input.CurrentImageUrl)
		return ad.Ad{}, nil
	})

	_, err := svc.Store(context.Background(), dto)
	assert.NoError(t, err)
}

func TestService_Store_WithImage_Error(t *testing.T) {
	ctrl, svc, _, _, mockStorage := setup(t)
	defer ctrl.Finish()

	image := &local.UploadedFile{Name: "bad.jpg"}
	dto := ad.StoreDTO{Image: image}

	mockStorage.EXPECT().Save(image).Return("", errors.New("upload error"))

	_, err := svc.Store(context.Background(), dto)
	assert.Error(t, err)
}

func TestService_Store_WithoutImage(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	dto := ad.StoreDTO{Image: nil}

	mockRepo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(ad.Ad{}, nil)

	_, err := svc.Store(context.Background(), dto)
	assert.NoError(t, err)
}

func TestService_Update_WithImage_Success(t *testing.T) {
	ctrl, svc, mockRepo, _, mockStorage := setup(t)
	defer ctrl.Finish()

	image := &local.UploadedFile{Name: "update.jpg"}
	dto := ad.UpdateDTO{Image: image}

	mockStorage.EXPECT().Save(image).Return("new/path.jpg", nil)
	mockRepo.EXPECT().Update(gomock.Any(), int64(1), gomock.Any()).DoAndReturn(func(ctx context.Context, id int64, input ad.UpdateDTO) error {
		assert.Equal(t, "new/path.jpg", input.CurrentImageUrl)
		return nil
	})

	err := svc.Update(context.Background(), 1, dto)
	assert.NoError(t, err)
}

func TestService_Update_ImageSaveError(t *testing.T) {
	ctrl, svc, _, _, mockStorage := setup(t)
	defer ctrl.Finish()

	image := &local.UploadedFile{Name: "bad.jpg"}
	dto := ad.UpdateDTO{Image: image}

	mockStorage.EXPECT().Save(image).Return("", errors.New("failed"))

	err := svc.Update(context.Background(), 1, dto)
	assert.Error(t, err)
}

func TestService_Update_WithoutImage(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	dto := ad.UpdateDTO{}

	mockRepo.EXPECT().Update(gomock.Any(), int64(1), dto).Return(nil)

	err := svc.Update(context.Background(), 1, dto)
	assert.NoError(t, err)
}

func TestService_Update_RepoError(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	dto := ad.UpdateDTO{}

	mockRepo.EXPECT().Update(gomock.Any(), int64(1), dto).Return(errors.New("db error"))

	err := svc.Update(context.Background(), 1, dto)
	assert.Error(t, err)
}

func TestService_List_Success(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	dto := ad.ListDTO{}
	expectedAds := []ad.Ad{{}}
	expectedCount := uint64(1)

	mockRepo.EXPECT().List(gomock.Any(), dto).Return(expectedAds, expectedCount, nil)

	ads, count, err := svc.List(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, expectedAds, ads)
	assert.Equal(t, expectedCount, count)
}

func TestService_List_Error(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	dto := ad.ListDTO{}

	mockRepo.EXPECT().List(gomock.Any(), dto).Return(nil, uint64(0), errors.New("list error"))

	_, _, err := svc.List(context.Background(), dto)
	assert.Error(t, err)
}

func TestService_Delete_Success(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)
	assert.NoError(t, err)
}

func TestService_Delete_Error(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().Delete(gomock.Any(), int64(1)).Return(errors.New("delete error"))

	err := svc.Delete(context.Background(), 1)
	assert.Error(t, err)
}

func TestService_HandleFavorite_Success(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().HandleFavorite(gomock.Any(), int64(1), int64(2)).Return(nil)

	err := svc.HandleFavorite(context.Background(), 1, 2)
	assert.NoError(t, err)
}

func TestService_HandleFavorite_Error(t *testing.T) {
	ctrl, svc, mockRepo, _, _ := setup(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().HandleFavorite(gomock.Any(), int64(1), int64(2)).Return(errors.New("favorite error"))

	err := svc.HandleFavorite(context.Background(), 1, 2)
	assert.Error(t, err)
}

func TestService_GetTokenInfo_Success(t *testing.T) {
	ctrl, svc, _, mockNftRepo, _ := setup(t)
	defer ctrl.Finish()

	expectedNFT := nft.NFT{TokenId: 123}
	mockNftRepo.EXPECT().GetNftInfo(gomock.Any(), big.NewInt(123)).Return(expectedNFT, nil)

	result, err := svc.GetTokenInfo(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, expectedNFT, result)
}

func TestService_GetTokenInfo_Error(t *testing.T) {
	ctrl, svc, _, mockNftRepo, _ := setup(t)
	defer ctrl.Finish()

	mockNftRepo.EXPECT().GetNftInfo(gomock.Any(), big.NewInt(123)).Return(nft.NFT{}, errors.New("nft error"))

	_, err := svc.GetTokenInfo(context.Background(), 123)
	assert.Error(t, err)
}
