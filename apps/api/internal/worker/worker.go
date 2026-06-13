package worker

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
	"github.com/ticketing/api/internal/service"
)

// StartExpiryJobs runs background jobs:
//   - Expire pending bookings (release stock) every 1 minute
//   - Expire pending payments (cancel associated bookings) every 1 minute
//   - Advance queue users to booking sessions every 10 seconds
//   - Expire abandoned queue tokens every 30 seconds
func StartExpiryJobs(ctx context.Context, db *gorm.DB, rdb *redis.Client, cfg *Config) {
	if rdb == nil {
		log.Println("worker: Redis not available, skipping queue jobs")
	} else {
		warSvc := service.NewWarQueue(rdb, cfg.QueueSessionTTLMinutes, cfg.WarRateLimitPerMin)
		qtRepo := repository.NewQueueTokenRepository(db)
		go runQueueAdvanceJob(ctx, warSvc)
		go runQueueCleanupJob(ctx, warSvc, qtRepo)
	}
	go runBookingExpiryJob(ctx, db)
	go runPaymentExpiryJob(ctx, db)
}

type Config struct {
	QueueSessionTTLMinutes int
	WarRateLimitPerMin     int
}

func runBookingExpiryJob(ctx context.Context, db *gorm.DB) {
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			expireBookings(ctx, bookingRepo, catRepo)
		}
	}
}

func expireBookings(ctx context.Context, bookingRepo *repository.BookingRepository, catRepo *repository.TicketCategoryRepository) {
	expired, err := bookingRepo.ListExpired(ctx)
	if err != nil {
		log.Printf("worker: list expired bookings: %v", err)
		return
	}

	for _, b := range expired {
		if err := bookingRepo.MarkExpired(ctx, b.ID); err != nil {
			log.Printf("worker: mark booking %d expired: %v", b.ID, err)
			continue
		}
		for _, item := range b.Items {
			if err := catRepo.ReleaseStock(ctx, item.TicketCategoryID, item.Quantity); err != nil {
				log.Printf("worker: release stock booking %d cat %d: %v", b.ID, item.TicketCategoryID, err)
			}
		}
		log.Printf("worker: expired booking %s (id=%d)", b.BookingCode, b.ID)
	}

	if len(expired) > 0 {
		log.Printf("worker: expired %d pending booking(s)", len(expired))
	}
}

func runPaymentExpiryJob(ctx context.Context, db *gorm.DB) {
	paymentRepo := repository.NewPaymentRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	catRepo := repository.NewTicketCategoryRepository(db)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			expirePayments(ctx, paymentRepo, bookingRepo, catRepo)
		}
	}
}

func expirePayments(ctx context.Context, paymentRepo *repository.PaymentRepository, bookingRepo *repository.BookingRepository, catRepo *repository.TicketCategoryRepository) {
	expired, err := paymentRepo.ListExpired(ctx)
	if err != nil {
		log.Printf("worker: list expired payments: %v", err)
		return
	}

	for _, p := range expired {
		if err := paymentRepo.UpdateStatus(ctx, p.ID, model.PaymentStatusExpired, nil, nil); err != nil {
			log.Printf("worker: mark payment %d expired: %v", p.ID, err)
			continue
		}

		booking, err := bookingRepo.GetByID(ctx, p.BookingID)
		if err != nil {
			log.Printf("worker: get booking %d for payment expiry: %v", p.BookingID, err)
			continue
		}

		if booking.Status == model.BookingStatusPending {
			if _, err := bookingRepo.CancelBooking(ctx, booking.ID, catRepo); err != nil {
				log.Printf("worker: cancel booking %d for expired payment: %v", booking.ID, err)
			} else {
				log.Printf("worker: cancelled booking %s (id=%d) due to expired payment", booking.BookingCode, booking.ID)
			}
		}
	}

	if len(expired) > 0 {
		log.Printf("worker: expired %d pending payment(s)", len(expired))
	}
}

// runQueueAdvanceJob checks every 10 seconds if the next user in queue should be advanced.
func runQueueAdvanceJob(ctx context.Context, war *service.WarQueue) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			advanceAllQueues(ctx, war)
		}
	}
}

func advanceAllQueues(ctx context.Context, war *service.WarQueue) {
	keys, err := war.ScanQueueKeys(ctx)
	if err != nil {
		log.Printf("worker: scan queue keys: %v", err)
		return
	}
	for _, eventID := range keys {
		n, err := war.ProcessQueue(ctx, eventID, 1)
		if err != nil {
			log.Printf("worker: advance queue event=%d: %v", eventID, err)
			continue
		}
		if n > 0 {
			log.Printf("worker: advanced %d user(s) from event %d", n, eventID)
		}
	}
}

// runQueueCleanupJob removes expired/abandoned queue entries every 30 seconds.
func runQueueCleanupJob(ctx context.Context, war *service.WarQueue, qtRepo *repository.QueueTokenRepository) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			keys, err := war.ScanQueueKeys(ctx)
			if err != nil {
				log.Printf("worker: scan queue keys for cleanup: %v", err)
				continue
			}
			for _, eventID := range keys {
				removed, err := war.CleanupExpired(ctx, eventID, 2*time.Minute)
				if err != nil {
					log.Printf("worker: cleanup queue event=%d: %v", eventID, err)
				}
				if removed > 0 {
					log.Printf("worker: removed %d expired entries from event %d", removed, eventID)
				}
			}

			expired, err := qtRepo.ExpireStaleTokens(ctx, 2)
			if err != nil {
				log.Printf("worker: expire stale db tokens: %v", err)
			}
			if expired > 0 {
				log.Printf("worker: expired %d stale queue tokens in db", expired)
			}
		}
	}
}
