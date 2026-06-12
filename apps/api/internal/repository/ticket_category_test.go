package repository_test

import (
	"context"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
)

// Integration test: requires a running Postgres.
// Set TEST_DATABASE_URL to a Postgres DSN.
//
// Example:
//   TEST_DATABASE_URL=postgres://ticketing:ticketing_secret@localhost:5432/ticketing_test?sslmode=disable \
//   go test ./internal/repository/...
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run integration tests")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	// Auto-migrate for test
	if err := db.AutoMigrate(
		&model.User{},
		&model.Event{},
		&model.TicketCategory{},
		&model.Booking{},
		&model.BookingItem{},
		&model.Payment{},
		&model.QueueToken{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func seedCategory(t *testing.T, db *gorm.DB, eventID uint64, stock int) uint64 {
	t.Helper()
	c := &model.TicketCategory{
		EventID:        eventID,
		Name:           "TestCat",
		Price:          100000,
		TotalStock:     stock,
		AvailableStock: stock,
		MaxPerUser:     4,
	}
	if err := db.Create(c).Error; err != nil {
		t.Fatalf("seed category: %v", err)
	}
	return c.ID
}

func TestReserveStock_Basic(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTicketCategoryRepository(db)

	event := &model.Event{Title: "T", Venue: "V", StartDate: time.Now(), EndDate: time.Now().Add(time.Hour), Status: model.EventStatusDraft}
	if err := db.Create(event).Error; err != nil {
		t.Fatalf("seed event: %v", err)
	}
	catID := seedCategory(t, db, event.ID, 10)

	// Reserve 3
	updated, err := repo.ReserveStock(context.Background(), catID, 3)
	if err != nil {
		t.Fatalf("reserve: %v", err)
	}
	if updated.AvailableStock != 7 {
		t.Errorf("avail: want 7, got %d", updated.AvailableStock)
	}

	// Reserve more than available
	if _, err := repo.ReserveStock(context.Background(), catID, 100); !errors.Is(err, repository.ErrInsufficientStock) {
		t.Errorf("want ErrInsufficientStock, got %v", err)
	}
}

func TestReleaseStock(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTicketCategoryRepository(db)

	event := &model.Event{Title: "T", Venue: "V", StartDate: time.Now(), EndDate: time.Now().Add(time.Hour), Status: model.EventStatusDraft}
	db.Create(event)
	catID := seedCategory(t, db, event.ID, 10)

	repo.ReserveStock(context.Background(), catID, 5)
	if err := repo.ReleaseStock(context.Background(), catID, 3); err != nil {
		t.Fatalf("release: %v", err)
	}
	cur, _ := repo.FindByID(context.Background(), catID)
	if cur.AvailableStock != 8 {
		t.Errorf("avail after release: want 8, got %d", cur.AvailableStock)
	}
}

// TestReserveStock_Concurrent simulates N goroutines reserving the last
// ticket; only 1 should succeed.
func TestReserveStock_Concurrent(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTicketCategoryRepository(db)

	event := &model.Event{Title: "T", Venue: "V", StartDate: time.Now(), EndDate: time.Now().Add(time.Hour), Status: model.EventStatusDraft}
	db.Create(event)
	catID := seedCategory(t, db, event.ID, 5) // only 5 tickets

	const goroutines = 1000
	var success, soldOut int64
	var wg sync.WaitGroup
	wg.Add(goroutines)

	ctx := context.Background()
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := repo.ReserveStock(ctx, catID, 1)
			switch {
			case err == nil:
				atomic.AddInt64(&success, 1)
			case errors.Is(err, repository.ErrInsufficientStock):
				atomic.AddInt64(&soldOut, 1)
			default:
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	if success != 5 {
		t.Errorf("expected exactly 5 successful reservations, got %d", success)
	}
	if soldOut != goroutines-5 {
		t.Errorf("expected %d soldOut, got %d", goroutines-5, soldOut)
	}

	// Verify final stock is 0
	cur, _ := repo.FindByID(ctx, catID)
	if cur.AvailableStock != 0 {
		t.Errorf("expected 0 stock remaining, got %d", cur.AvailableStock)
	}
}
