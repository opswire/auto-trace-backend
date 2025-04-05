package http

type BasicResponseDTO struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}
