package ad

import (
	"car-sell-buy-system/internal/ads-service/domain/ad"
	"car-sell-buy-system/pkg/storage/local"
	"mime/multipart"
)

type StoreRequest struct {
	Title         string                `form:"title"`
	Description   string                `form:"description"`
	Price         float64               `form:"price"`
	Vin           string                `form:"vin"`
	Brand         string                `form:"brand"`
	Model         string                `form:"model"`
	YearOfRelease int64                 `form:"year_of_release"`
	Image         *multipart.FileHeader `form:"image" binding:"required"`
}

func (r StoreRequest) ToDTO() (ad.StoreDTO, error) {
	var image *local.UploadedFile
	if r.Image != nil {
		img, err := local.ConvertUploadedFile(r.Image)
		if err == nil {
			image = img
		}
	}

	return ad.StoreDTO{
		Title:         r.Title,
		Description:   r.Description,
		Price:         r.Price,
		Vin:           r.Vin,
		Brand:         r.Brand,
		Model:         r.Model,
		YearOfRelease: r.YearOfRelease,
		Image:         image,
	}, nil
}

type UpdateRequest struct {
	Title           string                `form:"title" binding:"required"`
	Description     string                `form:"description" binding:"required"`
	Price           float64               `form:"price" binding:"required"`
	Vin             string                `form:"vin" binding:"required"`
	Brand           string                `form:"brand" binding:"required"`
	Model           string                `form:"model" binding:"required"`
	YearOfRelease   int64                 `form:"year_of_release" binding:"required"`
	Image           *multipart.FileHeader `form:"image"`
	CurrentImageUrl string                `form:"image_url" binding:"required"`
}

func (r UpdateRequest) ToDTO() (ad.UpdateDTO, error) {
	var image *local.UploadedFile
	if r.Image != nil {
		img, err := local.ConvertUploadedFile(r.Image)
		if err != nil {
			return ad.UpdateDTO{}, err
		}
		image = img
	}

	return ad.UpdateDTO{
		Title:           r.Title,
		Description:     r.Description,
		Price:           r.Price,
		Vin:             r.Vin,
		Brand:           r.Brand,
		Model:           r.Model,
		YearOfRelease:   r.YearOfRelease,
		Image:           image,
		CurrentImageUrl: r.CurrentImageUrl,
	}, nil
}

type HandleFavoriteRequest struct {
	AdId int64 `json:"ad_id"`
}

func (r HandleFavoriteRequest) ToDTO() ad.HandleFavoriteDTO {
	return ad.HandleFavoriteDTO{
		AdId: r.AdId,
	}
}
