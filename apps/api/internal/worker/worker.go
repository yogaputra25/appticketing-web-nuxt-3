package worker

import (
	"context"
	"gorm.io/gorm"
)

// StartExpiryJobs runs background jobs:
//   - Expire pending bookings (release stock) every 1 minute
//   - Expire queue tokens (>2 menit idle) every 30 detik
// Implementations are added in section 7 & 6.
func StartExpiryJobs(ctx context.Context, db *gorm.DB) {
	go runBookingExpiryJob(ctx, db)
	// go runQueueCleanupJob(ctx, db)
}

func runBookingExpiryJob(ctx context.Context, db *gorm.DB) {
	// Placeholder — implemented in section 7.7
	<-ctx.Done()
}
