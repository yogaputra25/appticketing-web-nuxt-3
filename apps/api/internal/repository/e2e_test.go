package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
	"github.com/ticketing/api/internal/service"
)

func setupTestRedisE2E(t *testing.T) *redis.Client {
	t.Helper()
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		t.Skip("set TEST_REDIS_URL to run integration tests")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatalf("parse redis url: %v", err)
	}
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("ping redis: %v", err)
	}
	return rdb
}

func cleanupQueueE2E(t *testing.T, rdb *redis.Client) {
	t.Helper()
	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, "war:*", 100).Iterator()
	for iter.Next(ctx) {
		rdb.Del(ctx, iter.Val())
	}
}

func TestE2EFullFlow(t *testing.T) {
	db := setupTestDB(t)
	rdb := setupTestRedisE2E(t)
	defer cleanupQueueE2E(t, rdb)

	ctx := context.Background()

	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	warSvc := service.NewWarQueue(rdb, 5, 10)

	// 1. Register user
	passwordHash, err := model.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := &model.User{
		Email:        "e2e@test.com",
		PasswordHash: passwordHash,
		FullName:     "E2E User",
		Role:         "user",
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	t.Logf("1. User created: id=%d", user.ID)

	// 2. Create event with category (admin action)
	desc := "E2E Test Event"
	event := &model.Event{
		Title:       "E2E Test Concert",
		Description: &desc,
		Venue:       "Test Venue",
		StartDate:   time.Now().Add(30 * 24 * time.Hour),
		EndDate:     time.Now().Add(31 * 24 * time.Hour),
		Status:      model.EventStatusPublished,
	}
	if err := eventRepo.Create(ctx, event); err != nil {
		t.Fatalf("create event: %v", err)
	}
	t.Logf("2. Event created: id=%d", event.ID)

	cat := &model.TicketCategory{
		EventID:        event.ID,
		Name:           "E2E Category",
		Price:          150000,
		TotalStock:     100,
		AvailableStock: 100,
		MaxPerUser:     4,
	}
	if err := catRepo.Create(ctx, cat); err != nil {
		t.Fatalf("create category: %v", err)
	}
	t.Logf("3. Category created: id=%d, price=%.0f", cat.ID, cat.Price)

	// 4. Join war queue
	_, token, err := warSvc.JoinQueue(ctx, user.ID, event.ID)
	if err != nil {
		t.Fatalf("join queue: %v", err)
	}
	t.Logf("4. Queue joined: token=%s", token)

	// 5. Check queue status
	pos, _, _, _, err := warSvc.GetQueueStatus(ctx, user.ID, event.ID)
	if err != nil {
		t.Fatalf("get queue status: %v", err)
	}
	t.Logf("5. Queue status: position=%d", pos)

	// 6. Process queue (advance user to booking session)
	n, err := warSvc.ProcessQueue(ctx, event.ID, 1)
	if err != nil {
		t.Fatalf("process queue: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 user advanced, got %d", n)
	}
	t.Log("6. Queue advanced")

	// 7. Issue booking session
	sessionToken, err := warSvc.IssueBookingSession(ctx, user.ID, event.ID)
	if err != nil {
		t.Fatalf("issue booking session: %v", err)
	}
	t.Logf("7. Session issued: token=%s", sessionToken)

	// 8. Validate session and create booking
	uid, eid, err := warSvc.ValidateBookingSession(ctx, sessionToken)
	if err != nil {
		t.Fatalf("validate session: err=%v", err)
	}
	if uid != user.ID || eid != event.ID {
		t.Fatalf("validate session: expected user=%d event=%d, got user=%d event=%d", user.ID, event.ID, uid, eid)
	}

	booking, err := bookingRepo.Create(ctx, repository.CreateBookingInput{
		UserID:  user.ID,
		EventID: event.ID,
		Items: []repository.BookingItemInput{
			{CategoryID: cat.ID, Quantity: 2, UnitPrice: cat.Price},
		},
	}, catRepo, 10)
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}
	t.Logf("8. Booking created: code=%s, total=%.0f", booking.BookingCode, booking.TotalAmount)

	// 9. Check stock was reserved
	updatedCat, err := catRepo.FindByID(ctx, cat.ID)
	if err != nil {
		t.Fatalf("find category: %v", err)
	}
	if updatedCat.AvailableStock != 98 {
		t.Errorf("expected available 98, got %d", updatedCat.AvailableStock)
	}
	t.Logf("9. Stock reserved: available=%d", updatedCat.AvailableStock)

	// 10. Create payment
	payment := &model.Payment{
		PaymentCode:   "E2E-PAY-001",
		BookingID:     booking.ID,
		UserID:        user.ID,
		Amount:        booking.TotalAmount,
		Status:        model.PaymentStatusPending,
		PaymentMethod: "simulation",
		ExpiredAt:     timePtr(time.Now().Add(30 * time.Minute)),
	}
	if err := paymentRepo.Create(ctx, payment); err != nil {
		t.Fatalf("create payment: %v", err)
	}
	t.Logf("10. Payment created: id=%d", payment.ID)

	// 11. Simulate payment success
	now := time.Now()
	if err := paymentRepo.UpdateStatus(ctx, payment.ID, model.PaymentStatusSuccess, &now, map[string]string{
		"simulated_at": now.Format(time.RFC3339),
		"result":       "approved",
	}); err != nil {
		t.Fatalf("simulate payment success: %v", err)
	}

	// 12. Update booking to paid and add e-ticket codes
	if err := bookingRepo.UpdateStatus(ctx, booking.ID, model.BookingStatusPaid); err != nil {
		t.Fatalf("update booking to paid: %v", err)
	}
	eTicketCodes := []string{"TCK-E2E-001", "TCK-E2E-002"}
	if err := bookingRepo.UpdateETicketCodes(ctx, booking.ID, eTicketCodes); err != nil {
		t.Fatalf("update e-ticket codes: %v", err)
	}
	t.Logf("11-12. Payment + e-ticket: codes=%v", eTicketCodes)

	// 13. Verify final state
	finalBooking, err := bookingRepo.GetByID(ctx, booking.ID)
	if err != nil {
		t.Fatalf("get final booking: %v", err)
	}
	if finalBooking.Status != model.BookingStatusPaid {
		t.Errorf("expected booking paid, got %s", finalBooking.Status)
	}
	finalPayment, err := paymentRepo.GetByID(ctx, payment.ID)
	if err != nil {
		t.Fatalf("get final payment: %v", err)
	}
	if finalPayment.Status != model.PaymentStatusSuccess {
		t.Errorf("expected payment success, got %s", finalPayment.Status)
	}
	t.Log("13. Final state verified: booking=paid, payment=success")
}
