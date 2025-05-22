package nft

import (
	"car-sell-buy-system/internal/nft-service/domain/nft"
)

type StoreNftRequest struct {
	Vin         string `json:"vin" binding:"required"`
	MetadataUrl string `json:"metadata_url"`
}

func (r StoreNftRequest) ToDTO() nft.StoreNftDTO {
	return nft.StoreNftDTO{
		Vin:         r.Vin,
		MetadataUrl: r.MetadataUrl,
	}
}

type AddServiceRecordRequest struct {
	Description string `json:"description" binding:"required"`
	Company     string `json:"company" binding:"required"`
	Signature   string `json:"signature" binding:"required"`
}

func (r AddServiceRecordRequest) ToDTO() nft.AddServiceRecordDTO {
	return nft.AddServiceRecordDTO{
		Description: r.Description,
		Company:     r.Company,
		Signature:   r.Signature,
	}
}
