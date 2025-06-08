package main

import (
	"bytes"
	"car-sell-buy-system/config"
	"car-sell-buy-system/internal/payments-service/domain/payment"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"html/template"
	"log"
	"net/smtp"
)

const paymentReceiptHTML = `
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <title>Чек об оплате</title>
  <style>
    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      background-color: #f4f4f4;
      padding: 40px;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background-color: #ffffff;
      border-radius: 10px;
      padding: 30px;
      box-shadow: 0 5px 15px rgba(0,0,0,0.05);
    }
    h1 {
      text-align: center;
      color: #2c3e50;
    }
    .info {
      margin: 20px 0;
    }
    .info p {
      font-size: 16px;
      margin: 8px 0;
      color: #333333;
    }
    .footer {
      margin-top: 30px;
      text-align: center;
      font-size: 14px;
      color: #999999;
    }
    .amount {
      font-size: 24px;
      font-weight: bold;
      color: #27ae60;
      margin-top: 20px;
      text-align: center;
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Чек об оплате</h1>

    <div class="info">
      <p><strong>Номер платежа:</strong> {{.PaymentID}}</p>
      <p><strong>Пользователь:</strong> {{.UserEmail}}</p>
      <p><strong>Объявление:</strong> {{.AdTitle}}</p>
    </div>

    <div class="amount">
      Сумма оплаты: {{.Amount}} ₽
    </div>

    <div class="footer">
      Спасибо за оплату!<br>
      Это письмо подтверждает успешное завершение транзакции.
    </div>
  </div>
</body>
</html>
`

func main() {
	cfg := config.NewConfig()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Kafka.Dsn},
		Topic:   "payments",
		GroupID: "email-sender",
	})

	log.Println("Consumer started...")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		var event payment.ConfirmedEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Println("JSON error:", err)
			continue
		}

		err = SendReceiptEmail(event, cfg)
		if err != nil {
			log.Println("SMTP error:", err)
		} else {
			log.Printf("Письмо отправлено на %s", event.UserEmail)
		}
	}
}

func SendReceiptEmail(event payment.ConfirmedEvent, cfg *config.Config) error {
	tmpl, err := template.New("receipt").Parse(paymentReceiptHTML)
	if err != nil {
		return fmt.Errorf("template error: %w", err)
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, event)
	if err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	return smtp.SendMail(fmt.Sprintf("%s:%s", cfg.Smtp.Host, cfg.Smtp.Port), nil,
		cfg.Smtp.From, []string{event.UserEmail},
		[]byte("Subject: Чек об оплате\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"+
			body.String()))
}
