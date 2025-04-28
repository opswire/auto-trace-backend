package yookassa

type (
	CreatePaymentRequest struct {
		Amount            AmountResponse           `json:"amount"`
		PaymentMethodData PaymentMethodDataRequest `json:"payment_method_data"`
		Confirmation      ConfirmationResponse     `json:"confirmation"`
		Description       string                   `json:"description"`
	}

	PaymentMethodDataRequest struct {
		Type string `json:"type"`
	}
)
