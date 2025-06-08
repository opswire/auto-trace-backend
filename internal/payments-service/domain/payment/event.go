package payment

type ConfirmedEvent struct {
	PaymentID string `json:"payment_id"`
	UserEmail string `json:"user_email"`
	Amount    int    `json:"amount"`
	AdTitle   string `json:"ad_title"`
}
