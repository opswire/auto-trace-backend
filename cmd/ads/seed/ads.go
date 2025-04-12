package main

import (
	"car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/pkg/postgres"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jaswdr/faker"
)

func SeedAds(ctx context.Context, pg *postgres.Postgres) error {
	fake := faker.New()
	rand.Seed(time.Now().UnixNano())

	// Список автомобильных брендов и моделей
	brands := []string{"Toyota", "Honda", "Ford", "BMW", "Mercedes", "Audi", "Tesla"}
	models := map[string][]string{
		"Toyota":   {"Camry", "Corolla", "RAV4"},
		"Honda":    {"Civic", "Accord", "CR-V"},
		"Ford":     {"Focus", "Mustang", "Explorer"},
		"BMW":      {"X5", "3 Series", "i8"},
		"Mercedes": {"C-Class", "E-Class", "GLA"},
		"Audi":     {"A4", "Q7", "TT"},
		"Tesla":    {"Model S", "Model 3", "Model X"},
	}

	// SQL для вставки
	query := `INSERT INTO ads (
		title, 
		description, 
		price, 
		vin, 
		is_token_minted, 
		brand, 
		model, 
		year_of_release
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	// Генерация 100 тестовых записей
	for i := 0; i < 100; i++ {
		brand := brands[rand.Intn(len(brands))]
		model := models[brand][rand.Intn(len(models[brand]))]

		ad := entity.Ad{
			Title:         fake.Lorem().Word(),
			Description:   fake.Lorem().Paragraph(3),
			Price:         rand.Float64() * 100000,
			Vin:           fake.Car().Plate(),
			IsTokenMinted: false, // 30% true
			Brand:         brand,
			Model:         model,
			YearOfRelease: int64(2000 + rand.Intn(24)), // 2000-2023
		}

		_, err := pg.Pool.Exec(ctx, query,
			ad.Title,
			ad.Description,
			ad.Price,
			ad.Vin,
			ad.IsTokenMinted,
			ad.Brand,
			ad.Model,
			ad.YearOfRelease,
		)

		if err != nil {
			return err
		}
	}

	fmt.Println("Successfully seeded 100 ads")

	return nil
}
