package repository_test

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
)

func seedUser(t *testing.T, db *gorm.DB) uint64 {
	t.Helper()
	hash, _ := model.HashPassword("password123")
	u := &model.User{
		Email:        "testuser@example.com",
		PasswordHash: hash,
		FullName:     "Test User",
		Role:         "user",
	}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u.ID
}

func seedBooking(t *testing.T, db *gorm.DB, userID, eventID uint64) uint64 {
	t.Helper()
	b := &model.Booking{
		BookingCode: "TEST-BK-" + t.Name(),
		UserID:      userID,
		EventID:     eventID,
		TotalAmount: 100000,
		Status:      model.BookingStatusPending,
		ExpiresAt:   timePtr(time.Now().Add(30 * time.Minute)),
	}
	if err := db.Create(b).Error; err != nil {
		t.Fatalf("seed booking: %v", err)
	}
	return b.ID
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestCreatePayment(t *testing.T) {
	db := setupTestDB(t)
	paymentRepo := repository.NewPaymentRepository(db)
	ctx := context.Background()

	userID := seedUser(t, db)
	eventID := seedEvent(t, db)
	bookingID := seedBooking(t, db, userID, eventID)

	now := time.Now()
	payment := &model.Payment{
		PaymentCode:   "PAY-TEST-001",
		BookingID:     bookingID,
		UserID:        userID,
		Amount:        100000,
		Status:        model.PaymentStatusPending,
		PaymentMethod: "simulation",
		ExpiredAt:     timePtr(now.Add(30 * time.Minute)),
	}

	if err := paymentRepo.Create(ctx, payment); err != nil {
		t.Fatalf("create payment: %v", err)
	}
	if payment.ID == 0 {
		t.Error("expected payment ID after create")
	}
}

func TestGetPaymentByID(t *testing.T) {
	db := setupTestDB(t)
	paymentRepo := repository.NewPaymentRepository(db)
	ctx := context.Background()

	userID := seedUser(t, db)
	eventID := seedEvent(t, db)
	bookingID := seedBooking(t, db, userID, eventID)

	payment := &model.Payment{
		PaymentCode:   "PAY-TEST-002",
		BookingID:     bookingID,
		UserID:        userID,
		Amount:        100000,
		Status:        model.PaymentStatusPending,
		PaymentMethod: "simulation",
	}

	if err := paymentRepo.Create(ctx, payment); err != nil {
		t.Fatalf("create payment: %v", err)
	}

	got, err := paymentRepo.GetByID(ctx, payment.ID)
	if err != nil {
		t.Fatalf("get payment: %v", err)
	}
	if got.PaymentCode != "PAY-TEST-002" {
		t.Errorf("expected code PAY-TEST-002, got %s", got.PaymentCode)
	}
	if got.Status != model.PaymentStatusPending {
		t.Errorf("expected status pending, got %s", got.Status)
	}
}

func TestUpdatePaymentStatus(t *testing.T) {
	db := setupTestDB(t)
	paymentRepo := repository.NewPaymentRepository(db)
	ctx := context.Background()

	userID := seedUser(t, db)
	eventID := seedEvent(t, db)
	bookingID := seedBooking(t, db, userID, eventID)

	payment := &model.Payment{
		PaymentCode:   "PAY-TEST-003",
		BookingID:     bookingID,
		UserID:        userID,
		Amount:        100000,
		Status:        model.PaymentStatusPending,
		PaymentMethod: "simulation",
	}

	if err := paymentRepo.Create(ctx, payment); err != nil {
		t.Fatalf("create payment: %v", err)
	}

	now := time.Now()
	if err := paymentRepo.UpdateStatus(ctx, payment.ID, model.PaymentStatusSuccess, &now, nil); err != nil {
		t.Fatalf("update payment status: %v", err)
	}

	got, err := paymentRepo.GetByID(ctx, payment.ID)
	if err != nil {
		t.Fatalf("get payment: %v", err)
	}
	if got.Status != model.PaymentStatusSuccess {
		t.Errorf("expected status success, got %s", got.Status)
	}
	if got.PaidAt == nil {
		t.Error("expected paid_at to be set on success")
	}
}

func TestPaymentLifecycle(t *testing.T) {
	db := setupTestDB(t)
	paymentRepo := repository.NewPaymentRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	userID := seedUser(t, db)
	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 10)

	cat, err := catRepo.FindByID(ctx, catID)
	if err != nil {
		t.Fatalf("find category: %v", err)
	}

	booking, err := bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  userID,
		EventID: eventID,
		Items: []repository.BookingItemInput{
			{CategoryID: catID, Quantity: 2, UnitPrice: cat.Price},
		},
	}, catRepo, 10)
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	now := time.Now()
	payment := &model.Payment{
		PaymentCode:   "PAY-TEST-LIFECYCLE",
		BookingID:     booking.ID,
		UserID:        userID,
		Amount:        booking.TotalAmount,
		Status:        model.PaymentStatusPending,
		PaymentMethod: "simulation",
		ExpiredAt:     timePtr(now.Add(30 * time.Minute)),
	}

	if err := paymentRepo.Create(ctx, payment); err != nil {
		t.Fatalf("create payment: %v", err)
	}

	if err := paymentRepo.UpdateStatus(ctx, payment.ID, model.PaymentStatusSuccess, &now, map[string]string{
		"simulated_at": now.Format(time.RFC3339),
		"result":       "approved",
	}); err != nil {
		t.Fatalf("mark payment success: %v", err)
	}

	if err := bookingRepo.UpdateStatus(ctx, booking.ID, model.BookingStatusPaid); err != nil {
		t.Fatalf("update booking to paid: %v", err)
	}

	eTicketCodes := []string{"TCK-TEST-001", "TCK-TEST-002"}
	if err := bookingRepo.UpdateETicketCodes(ctx, booking.ID, eTicketCodes); err != nil {
		t.Fatalf("update e-ticket codes: %v", err)
	}

	updatedBooking, err := bookingRepo.GetByID(ctx, booking.ID)
	if err != nil {
		t.Fatalf("get booking: %v", err)
	}
	if updatedBooking.Status != model.BookingStatusPaid {
		t.Errorf("expected booking paid, got %s", updatedBooking.Status)
	}

	got, err := paymentRepo.GetByBookingID(ctx, booking.ID)
	if err != nil {
		t.Fatalf("get payment by booking: %v", err)
	}
	if got.Status != model.PaymentStatusSuccess {
		t.Errorf("expected payment success, got %s", got.Status)
	}
}

func TestCancelBookingOnPaymentExpiry(t *testing.T) {
	db := setupTestDB(t)
	paymentRepo := repository.NewPaymentRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	ctx := context.Background()

	userID := seedUser(t, db)
	eventID := seedEvent(t, db)
	catID := seedCategory(t, db, eventID, 10)

	cat, err := catRepo.FindByID(ctx, catID)
	if err != nil {
		t.Fatalf("find category: %v", err)
	}

	booking, err := bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  userID,
		EventID: eventID,
		Items: []repository.BookingItemInput{
			{CategoryID: catID, Quantity: 1, UnitPrice: cat.Price},
		},
	}, catRepo, 10)
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}

	payment := &model.Payment{
		PaymentCode:   "PAY-TEST-EXPIRY",
		BookingID:     booking.ID,
		UserID:        userID,
		Amount:        booking.TotalAmount,
		Status:        model.PaymentStatusPending,
		PaymentMethod: "simulation",
		ExpiredAt:     timePtr(time.Now().Add(-1 * time.Minute)),
	}

	if err := paymentRepo.Create(ctx, payment); err != nil {
		t.Fatalf("create payment: %v", err)
	}

	expiredPayments, err := paymentRepo.ListExpired(ctx)
	if err != nil {
		t.Fatalf("list expired: %v", err)
	}
	if len(expiredPayments) == 0 {
		t.Fatal("expected at least 1 expired payment")
	}

	found := false
	for _, p := range expiredPayments {
		if p.ID == payment.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected payment in expired list")
	}
}

func TestListPaymentsByUser(t *testing.T) {
	db := setupTestDB(t)
	paymentRepo := repository.NewPaymentRepository(db)
	ctx := context.Background()

	userID := seedUser(t, db)
	eventID := seedEvent(t, db)
	bookingID := seedBooking(t, db, userID, eventID)

	for i := 0; i < 3; i++ {
		payment := &model.Payment{
			PaymentCode:   "PAY-LIST-" + string(rune('A'+i)),
			BookingID:     bookingID,
			UserID:        userID,
			Amount:        100000,
			Status:        model.PaymentStatusPending,
			PaymentMethod: "simulation",
		}
		if err := paymentRepo.Create(ctx, payment); err != nil {
			t.Fatalf("create payment %d: %v", i, err)
		}
	}

	payments, total, err := paymentRepo.ListByUser(ctx, userID, 1, 10)
	if err != nil {
		t.Fatalf("list by user: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(payments) != 3 {
		t.Errorf("expected 3 payments, got %d", len(payments))
	}
}

func TestGetPaymentByBookingID(t *testing.T) {
	db := setupTestDB(t)
	paymentRepo := repository.NewPaymentRepository(db)
	ctx := context.Background()

	userID := seedUser(t, db)
	eventID := seedEvent(t, db)
	bookingID := seedBooking(t, db, userID, eventID)

	payment := &model.Payment{
		PaymentCode:   "PAY-BY-BK",
		BookingID:     bookingID,
		UserID:        userID,
		Amount:        100000,
		Status:        model.PaymentStatusPending,
		PaymentMethod: "simulation",
	}
	if err := paymentRepo.Create(ctx, payment); err != nil {
		t.Fatalf("create payment: %v", err)
	}

	got, err := paymentRepo.GetByBookingID(ctx, bookingID)
	if err != nil {
		t.Fatalf("get by booking: %v", err)
	}
	if got.ID != payment.ID {
		t.Errorf("expected payment id %d, got %d", payment.ID, got.ID)
	}
}
