package main

import (
	entity2 "car-sell-buy-system/internal/ads-service/entity"
	"car-sell-buy-system/internal/sso-service/entity"
	"car-sell-buy-system/pkg/postgres"
	"context"
	"fmt"
	"github.com/jaswdr/faker"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

func main() {
	ctx := context.Background()

	//cfg := config.NewConfig()

	//pg, err := postgres.New(cfg.Pg.URL)
	pg, err := postgres.New("postgres://user:pass@localhost:5454/postgres?sslmode=disable") // todo: refactor
	if err != nil {
		log.Fatal("Failed to connect DB.", err)
	}
	defer pg.Pool.Close()

	//tx, err := pg.Pool.Begin(ctx)
	//if err != nil {
	//	log.Fatal("Failed to start transaction.", err)
	//}

	if err = clearTables(ctx, pg); err != nil {
		fmt.Printf("Error clearing db: %v\n", err)
		//if err = tx.Rollback(ctx); err != nil {
		//	log.Fatal("Failed to rollback transaction.", err)
		//}
	}

	if err = seedUsers(ctx, pg); err != nil {
		fmt.Printf("Error seeding users: %v\n", err)
		//if err = tx.Rollback(ctx); err != nil {
		//	log.Fatal("Failed to rollback transaction.", err)
		//}
	}

	if err = seedAds(ctx, pg); err != nil {
		fmt.Printf("Error seeding ads: %v\n", err)
		//if err = tx.Rollback(ctx); err != nil {
		//	log.Fatal("Failed to rollback transaction.", err)
		//}
	}

	//if err = tx.Commit(ctx); err != nil {
	//	log.Fatal("Failed to commit transaction.", err)
	//}

	fmt.Println("Successfully seed!!!!!")
}

func clearTables(ctx context.Context, pg *postgres.Postgres) error {
	_, err := pg.Pool.Exec(ctx, "TRUNCATE TABLE users, ads RESTART IDENTITY CASCADE")
	return err
}

func seedUsers(ctx context.Context, pg *postgres.Postgres) error {
	users := []entity.User{
		{
			Email:    "user@example.com",
			Password: "user123",
			Role:     "user",
		},
		{
			Email:    "admin@example.com",
			Password: "admin123",
			Role:     "admin",
		},
		{
			Email:    "org@example.com",
			Password: "org123",
			Role:     "certified_organization",
		},
	}

	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3)`

	for _, u := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		_, err = pg.Pool.Exec(ctx, query, u.Email, string(hashedPassword), u.Role)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
	}

	fmt.Println("Seeded 3 base users:")
	fmt.Println("| Email               | Role                  |")
	fmt.Println("|---------------------|-----------------------|")
	fmt.Printf("| %-20s | %-21s |\n", "user@example.com", "user")
	fmt.Printf("| %-20s | %-21s |\n", "admin@example.com", "admin")
	fmt.Printf("| %-20s | %-21s |\n", "org@example.com", "certified_organization")
	return nil
}

func seedAds(ctx context.Context, pg *postgres.Postgres) error {
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

		ad := entity2.Ad{
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
