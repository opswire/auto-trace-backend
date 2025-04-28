package yookassa

import "time"

type (
	PaymentResponse struct {
		Id            string                `json:"id"`
		Status        string                `json:"status"`
		Paid          bool                  `json:"paid"`
		Amount        AmountResponse        `json:"amount"`
		Confirmation  ConfirmationResponse  `json:"confirmation"`
		Description   string                `json:"description"`
		UpdatedAt     time.Time             `json:"updated_at"`
		CreatedAt     time.Time             `json:"created_at"`
		ExpiresAt     time.Time             `json:"expires_at"`
		PaymentMethod PaymentMethodResponse `json:"payment_method"`
		Recipient     RecipientResponse     `json:"recipient"`
		Refundable    bool                  `json:"refundable"`
		Test          bool                  `json:"test"`
	}

	AmountResponse struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	}

	ConfirmationResponse struct {
		Type            string `json:"type"`
		ReturnUrl       string `json:"return_url"`
		ConfirmationUrl string `json:"confirmation_url"`
	}

	PaymentMethodResponse struct {
		PaymentMethodType string       `json:"type"`
		Id                string       `json:"id"`
		Saved             bool         `json:"saved"`
		Card              CardResponse `json:"card"`
		Title             string       `json:"title"`
	}

	CardResponse struct {
		First6        string              `json:"first6"`
		Last4         string              `json:"last4"`
		ExpiryMonth   string              `json:"expiry_month"`
		ExpiryYear    string              `json:"expiry_year"`
		CardType      string              `json:"card_type"`
		CardProduct   CardProductResponse `json:"card_product"`
		IssuerCountry string              `json:"issuer_country"`
		IssuerName    string              `json:"issuer_name"`
	}

	CardProductResponse struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	RecipientResponse struct {
		AccountId string `json:"account_id"`
		GatewayId string `json:"gateway_id"`
	}
)
