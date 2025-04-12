package handler

type BasicResponseDTO struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type BasicWithMetaResponseDTO struct {
	BasicResponseDTO
	Meta interface{} `json:"meta"`
}
