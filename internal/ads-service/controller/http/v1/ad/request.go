package ad

import (
	"car-sell-buy-system/internal/ads-service/domain/ad"
	"car-sell-buy-system/pkg/storage/local"
	"mime/multipart"
)

type StoreRequest struct {
	Title         string                `form:"title" binding:"required"`
	Description   string                `form:"description" binding:"required"`
	Price         float64               `form:"price" binding:"required"`
	Vin           string                `form:"vin" binding:"required"`
	Brand         string                `form:"brand" binding:"required"`
	Model         string                `form:"model" binding:"required"`
	YearOfRelease int64                 `form:"year_of_release"`
	Image         *multipart.FileHeader `form:"image" binding:"required"`
	Category      string                `form:"category" binding:"required"`
	RegNumber     string                `form:"reg_number" binding:"required"`
	Type          string                `form:"type" binding:"required"`
	Color         string                `form:"color" binding:"required"`
	Hp            string                `form:"hp" binding:"required"`
	FullWeight    string                `form:"full_weight" binding:"required"`
	SoloWeight    string                `form:"solo_weight" binding:"required"`
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
		Category:      r.Category,
		RegNumber:     r.RegNumber,
		Type:          r.Type,
		Color:         r.Color,
		Hp:            r.Hp,
		FullWeight:    r.FullWeight,
		SoloWeight:    r.SoloWeight,
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
	Category        string                `form:"category" binding:"required"`
	RegNumber       string                `form:"reg_number" binding:"required"`
	Type            string                `form:"type" binding:"required"`
	Color           string                `form:"color" binding:"required"`
	Hp              string                `form:"hp" binding:"required"`
	FullWeight      string                `form:"full_weight" binding:"required"`
	SoloWeight      string                `form:"solo_weight" binding:"required"`
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
		Category:        r.Category,
		RegNumber:       r.RegNumber,
		Type:            r.Type,
		Color:           r.Color,
		Hp:              r.Hp,
		FullWeight:      r.FullWeight,
		SoloWeight:      r.SoloWeight,
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
