package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
)

func seedEvent(t *testing.T, db *gorm.DB) uint64 {
	t.Helper()
	e := &model.Event{
		Title:     "Test Event",
		Venue:     "Test Venue",
		StartDate: time.Now().Add(24 * time.Hour),
		EndDate:   time.Now().Add(25 * time.Hour),
		Status:    model.EventStatusPublished,
	}
	if err := db.Create(e).Error; err != nil {
		t.Fatalf("seed event: %v", err)
	}
	return e.ID
}

func TestCreateBooking(t *testing.T) {
	db := setupTestDB(t)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 10)
	cat, _ := catRepo.FindByID(ctx, catID)

	booking, err := bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  1,
		EventID: eventID,
		Items: []repository.BookingItemInput{
			{CategoryID: catID, Quantity: 2, UnitPrice: cat.Price},
		},
	}, catRepo, 10)
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	if booking.BookingCode == "" {
		t.Error("expected non-empty booking code")
	}
	if booking.Status != model.BookingStatusPending {
		t.Errorf("expected pending status, got %s", booking.Status)
	}
	if booking.TotalAmount != 200000 {
		t.Errorf("expected total 200000, got %f", booking.TotalAmount)
	}
	if len(booking.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(booking.Items))
	}

	updated, _ := catRepo.FindByID(ctx, catID)
	if updated.AvailableStock != 8 {
		t.Errorf("expected stock 8 after reserve 2, got %d", updated.AvailableStock)
	}
}

func TestCancelPendingBooking(t *testing.T) {
	db := setupTestDB(t)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 10)
	cat, _ := catRepo.FindByID(ctx, catID)

	booking, err := bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  1,
		EventID: eventID,
		Items:   []repository.BookingItemInput{{CategoryID: catID, Quantity: 3, UnitPrice: cat.Price}},
	}, catRepo, 10)
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	cancelled, err := bookingRepo.CancelBooking(ctx, booking.ID, catRepo)
	if err != nil {
		t.Fatalf("cancel booking: %v", err)
	}
	if cancelled.Status != model.BookingStatusCancelled {
		t.Errorf("expected cancelled status, got %s", cancelled.Status)
	}

	updated, _ := catRepo.FindByID(ctx, catID)
	if updated.AvailableStock != 10 {
		t.Errorf("expected stock restored to 10, got %d", updated.AvailableStock)
	}
}

func TestCancelPaidBooking(t *testing.T) {
	db := setupTestDB(t)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 10)
	cat, _ := catRepo.FindByID(ctx, catID)

	booking, err := bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  1,
		EventID: eventID,
		Items:   []repository.BookingItemInput{{CategoryID: catID, Quantity: 1, UnitPrice: cat.Price}},
	}, catRepo, 10)
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	bookingRepo.UpdateStatus(ctx, booking.ID, model.BookingStatusPaid)

	_, err = bookingRepo.CancelBooking(ctx, booking.ID, catRepo)
	if !errors.Is(err, repository.ErrBookingNotPending) {
		t.Errorf("expected ErrBookingNotPending, got %v", err)
	}
}

func TestGetByUser(t *testing.T) {
	db := setupTestDB(t)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 20)
	cat, _ := catRepo.FindByID(ctx, catID)

	bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  10,
		EventID: eventID,
		Items:   []repository.BookingItemInput{{CategoryID: catID, Quantity: 1, UnitPrice: cat.Price}},
	}, catRepo, 10)
	bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  10,
		EventID: eventID,
		Items:   []repository.BookingItemInput{{CategoryID: catID, Quantity: 2, UnitPrice: cat.Price}},
	}, catRepo, 10)

	bookings, total, err := bookingRepo.GetByUser(ctx, 10, 1, 10)
	if err != nil {
		t.Fatalf("get by user: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(bookings) != 2 {
		t.Errorf("expected 2 bookings, got %d", len(bookings))
	}
}

func TestGetByIDOwnership(t *testing.T) {
	db := setupTestDB(t)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 10)
	cat, _ := catRepo.FindByID(ctx, catID)

	booking, err := bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  20,
		EventID: eventID,
		Items:   []repository.BookingItemInput{{CategoryID: catID, Quantity: 1, UnitPrice: cat.Price}},
	}, catRepo, 10)
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	fetched, err := bookingRepo.GetByID(ctx, booking.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if fetched.UserID != 20 {
		t.Errorf("expected user 20, got %d", fetched.UserID)
	}
	if len(fetched.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(fetched.Items))
	}
}

func TestInsufficientStock(t *testing.T) {
	db := setupTestDB(t)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 3)
	cat, _ := catRepo.FindByID(ctx, catID)

	_, err := bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  1,
		EventID: eventID,
		Items:   []repository.BookingItemInput{{CategoryID: catID, Quantity: 10, UnitPrice: cat.Price}},
	}, catRepo, 10)
	if !errors.Is(err, repository.ErrInsufficientStock) {
		t.Errorf("expected ErrInsufficientStock, got %v", err)
	}
}

func TestStockRestoredAfterCancel(t *testing.T) {
	db := setupTestDB(t)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 5)
	cat, _ := catRepo.FindByID(ctx, catID)

	booking, _ := bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  1,
		EventID: eventID,
		Items:   []repository.BookingItemInput{{CategoryID: catID, Quantity: 2, UnitPrice: cat.Price}},
	}, catRepo, 10)

	before, _ := catRepo.FindByID(ctx, catID)
	if before.AvailableStock != 3 {
		t.Fatalf("expected 3 after reservation, got %d", before.AvailableStock)
	}

	bookingRepo.CancelBooking(ctx, booking.ID, catRepo)

	after, _ := catRepo.FindByID(ctx, catID)
	if after.AvailableStock != 5 {
		t.Errorf("expected stock restored to 5, got %d", after.AvailableStock)
	}
}

func TestListExpired(t *testing.T) {
	db := setupTestDB(t)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 10)
	cat, _ := catRepo.FindByID(ctx, catID)

	expiredTime := time.Now().Add(-1 * time.Hour)
	booking := &model.Booking{
		BookingCode: "BK-EXPIRED",
		UserID:      1,
		EventID:     eventID,
		TotalAmount: 100000,
		Status:      model.BookingStatusPending,
		ExpiresAt:   &expiredTime,
		Items: []model.BookingItem{
			{TicketCategoryID: catID, Quantity: 2, UnitPrice: cat.Price, Subtotal: 100000},
		},
	}
	if err := db.Create(booking).Error; err != nil {
		t.Fatalf("seed expired booking: %v", err)
	}

	expired, err := bookingRepo.ListExpired(ctx)
	if err != nil {
		t.Fatalf("list expired: %v", err)
	}
	if len(expired) == 0 {
		t.Error("expected at least 1 expired booking")
	}
}
