package yookassa

import "time"

type (
	Payment struct {
		Id            string        `json:"id"`
		Status        string        `json:"status"`
		Paid          bool          `json:"paid"`
		Amount        Amount        `json:"amount"`
		Confirmation  Confirmation  `json:"confirmation"`
		Description   string        `json:"description"`
		UpdatedAt     time.Time     `json:"updated_at"`
		CreatedAt     time.Time     `json:"created_at"`
		ExpiresAt     time.Time     `json:"expires_at"`
		PaymentMethod PaymentMethod `json:"payment_method"`
		Recipient     Recipient     `json:"recipient"`
		Refundable    bool          `json:"refundable"`
		Test          bool          `json:"test"`
	}

	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	}

	Confirmation struct {
		Type            string `json:"type"`
		ReturnUrl       string `json:"return_url"`
		ConfirmationUrl string `json:"confirmation_url"`
	}

	PaymentMethod struct {
		PaymentMethodType string `json:"type"`
		Id                string `json:"id"`
		Saved             bool   `json:"saved"`
		Card              Card   `json:"card"`
		Title             string `json:"title"`
	}

	Card struct {
		First6        string      `json:"first6"`
		Last4         string      `json:"last4"`
		ExpiryMonth   string      `json:"expiry_month"`
		ExpiryYear    string      `json:"expiry_year"`
		CardType      string      `json:"card_type"`
		CardProduct   CardProduct `json:"card_product"`
		IssuerCountry string      `json:"issuer_country"`
		IssuerName    string      `json:"issuer_name"`
	}

	CardProduct struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	Recipient struct {
		AccountId string `json:"account_id"`
		GatewayId string `json:"gateway_id"`
	}
)
