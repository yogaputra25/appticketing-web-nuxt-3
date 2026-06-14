package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ticketing/api/internal/config"
	"github.com/ticketing/api/internal/database"
	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "database error: %v\n", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Event{},
		&model.TicketCategory{},
		&model.Booking{},
		&model.BookingItem{},
		&model.Payment{},
		&model.QueueToken{},
		&model.Ticket{},
	); err != nil {
		fmt.Fprintf(os.Stderr, "migrate error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("migration complete")

	ctx := context.Background()
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)

	if err := seedAdmin(ctx, userRepo); err != nil {
		fmt.Fprintf(os.Stderr, "seed admin error: %v\n", err)
		os.Exit(1)
	}

	if err := seedSampleEvents(ctx, eventRepo, catRepo); err != nil {
		fmt.Fprintf(os.Stderr, "seed events error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("seed complete")
}

func seedAdmin(ctx context.Context, userRepo *repository.UserRepository) error {
	existing, err := userRepo.FindByEmail(ctx, "admin@example.com")
	if err == nil && existing != nil {
		fmt.Println("admin user already exists, skipping")
		return nil
	}

	passwordHash, err := model.HashPassword("admin123")
	if err != nil {
		return err
	}

	u := &model.User{
		Email:        "admin@example.com",
		PasswordHash: passwordHash,
		FullName:     "System Admin",
		Role:         "admin",
	}

	if err := userRepo.Create(ctx, u); err != nil {
		return err
	}

	fmt.Printf("admin user created: %s (id=%d)\n", u.Email, u.ID)
	return nil
}

func seedSampleEvents(ctx context.Context, eventRepo *repository.EventRepository, catRepo *repository.TicketCategoryRepository) error {
	events, _, err := eventRepo.List(ctx, repository.EventFilter{Limit: 1})
	if err != nil {
		return err
	}
	if len(events) > 0 {
		fmt.Println("sample events already exist, skipping")
		return nil
	}

	now := time.Now()
	sampleEvents := []struct {
		title, venue, desc string
		startOffset, endOffset int
		categories []model.TicketCategory
	}{
		{
			title: "Music Festival 2026",
			venue: "Jakarta Convention Center",
			desc:  "Experience an unforgettable night of music and entertainment featuring top artists.",
			startOffset: 1,
			endOffset:   2,
			categories: []model.TicketCategory{
				{Name: "VIP", Price: 500000, TotalStock: 100, AvailableStock: 100, MaxPerUser: 2},
				{Name: "Regular", Price: 200000, TotalStock: 500, AvailableStock: 500, MaxPerUser: 4},
				{Name: "Economy", Price: 100000, TotalStock: 1000, AvailableStock: 1000, MaxPerUser: 4},
			},
		},
		{
			title: "Tech Conference 2026",
			venue: "ICE BSD City",
			desc:  "The premier technology conference featuring keynote speakers, workshops, and networking.",
			startOffset: 2,
			endOffset:   3,
			categories: []model.TicketCategory{
				{Name: "VIP Pass", Price: 1500000, TotalStock: 50, AvailableStock: 50, MaxPerUser: 2},
				{Name: "Standard Pass", Price: 750000, TotalStock: 300, AvailableStock: 300, MaxPerUser: 4},
			},
		},
		{
			title: "Sport Championship 2026",
			venue: "Gelora Bung Karno",
			desc:  "Witness the biggest sporting event of the year with top athletes competing for glory.",
			startOffset: 1,
			endOffset:   2,
			categories: []model.TicketCategory{
				{Name: "VVIP", Price: 2000000, TotalStock: 30, AvailableStock: 30, MaxPerUser: 2},
				{Name: "Tribune", Price: 500000, TotalStock: 500, AvailableStock: 500, MaxPerUser: 4},
				{Name: "Economy", Price: 150000, TotalStock: 2000, AvailableStock: 2000, MaxPerUser: 6},
			},
		},
	}

	for _, se := range sampleEvents {
		desc := se.desc
		event := &model.Event{
			Title:       se.title,
			Description: &desc,
			Venue:       se.venue,
			StartDate:   now.AddDate(0, se.startOffset, 0),
			EndDate:     now.AddDate(0, se.endOffset, 0),
			Status:      model.EventStatusPublished,
		}
		if err := eventRepo.Create(ctx, event); err != nil {
			return err
		}

		for i := range se.categories {
			se.categories[i].EventID = event.ID
			if err := catRepo.Create(ctx, &se.categories[i]); err != nil {
				return err
			}
		}

		fmt.Printf("sample event created: %s (id=%d) with %d categories\n", event.Title, event.ID, len(se.categories))
	}

	return nil
}
