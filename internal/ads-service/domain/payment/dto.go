package payment

type (
	CreatePaymentDto struct {
		AdId     int64
		UserId   int64
		TariffId int64
	}

	ProcessWebhookPaymentDto struct {
		TransactionId string
		Status        string
		Paid          bool
	}
)
