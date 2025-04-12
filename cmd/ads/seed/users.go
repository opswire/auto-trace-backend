package main

import (
	"car-sell-buy-system/internal/sso-service/entity"
	"car-sell-buy-system/pkg/postgres"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func SeedUsers(ctx context.Context, pg *postgres.Postgres) error {
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

	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3) 
			ON CONFLICT (email) DO NOTHING`

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
