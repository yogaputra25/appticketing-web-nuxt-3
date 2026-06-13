package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/ticketing/api/internal/model"
)

var (
	ErrBookingNotFound   = errors.New("booking not found")
	ErrBookingNotPending = errors.New("booking is not in pending payment state")
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

type CreateBookingInput struct {
	UserID  uint64
	EventID uint64
	Items   []BookingItemInput
}

type BookingItemInput struct {
	CategoryID uint64
	Quantity   int
	UnitPrice  float64
}

// Create creates a booking with items inside a transaction, reserving stock atomically.
func (r *BookingRepository) Create(ctx context.Context, input CreateBookingInput, catRepo *TicketCategoryRepository, ttlMinutes int) (*model.Booking, error) {
	if len(input.Items) == 0 {
		return nil, errors.New("at least one ticket item required")
	}

	var booking *model.Booking
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		itemMap := make(map[uint64]int)
		for _, it := range input.Items {
			itemMap[it.CategoryID] += it.Quantity
		}

		reserved, err := catRepo.ReserveStockMulti(ctx, itemMap)
		if err != nil {
			return err
		}

		code := generateBookingCode()
		expiresAt := time.Now().Add(time.Duration(ttlMinutes) * time.Minute)

		var totalAmount float64
		var items []model.BookingItem
		for _, it := range input.Items {
			if _, ok := reserved[it.CategoryID]; !ok {
				continue
			}
			subtotal := float64(it.Quantity) * it.UnitPrice
			totalAmount += subtotal
			items = append(items, model.BookingItem{
				TicketCategoryID: it.CategoryID,
				Quantity:         it.Quantity,
				UnitPrice:        it.UnitPrice,
				Subtotal:         subtotal,
			})
		}

		booking = &model.Booking{
			BookingCode: code,
			UserID:      input.UserID,
			EventID:     input.EventID,
			TotalAmount: totalAmount,
			Status:      model.BookingStatusPending,
			ExpiresAt:   &expiresAt,
			Items:       items,
		}

		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return booking, nil
}

// GetByUser returns paginated bookings for a user.
func (r *BookingRepository) GetByUser(ctx context.Context, userID uint64, page, limit int) ([]model.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	q := r.db.WithContext(ctx).Model(&model.Booking{}).Where("user_id = ?", userID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var bookings []model.Booking
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Preload("Event").Preload("Items").Find(&bookings).Error; err != nil {
		return nil, 0, err
	}
	return bookings, total, nil
}

// GetByID returns a booking by ID with items preloaded.
func (r *BookingRepository) GetByID(ctx context.Context, id uint64) (*model.Booking, error) {
	var b model.Booking
	if err := r.db.WithContext(ctx).Preload("Items").First(&b, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}
	return &b, nil
}

// GetByIDWithAssociations returns a booking by ID with items, user, and event preloaded.
func (r *BookingRepository) GetByIDWithAssociations(ctx context.Context, id uint64) (*model.Booking, error) {
	var b model.Booking
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("User").
		Preload("Event").
		First(&b, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}
	return &b, nil
}

// UpdateStatus changes the booking status.
func (r *BookingRepository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	res := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("id = ?", id).
		Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrBookingNotFound
	}
	return nil
}

// UpdateETicketCodes sets the e-ticket codes for a booking.
func (r *BookingRepository) UpdateETicketCodes(ctx context.Context, id uint64, codes []string) error {
	return r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("id = ?", id).
		Update("e_ticket_codes", model.JSONStringList(codes)).Error
}

// CancelBooking cancels a pending booking and releases stock.
func (r *BookingRepository) CancelBooking(ctx context.Context, id uint64, catRepo *TicketCategoryRepository) (*model.Booking, error) {
	var b *model.Booking
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var booking model.Booking
		if err := tx.Preload("Items").First(&booking, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBookingNotFound
			}
			return err
		}

		if booking.Status != model.BookingStatusPending && booking.Status != model.BookingStatusExpired {
			return ErrBookingNotPending
		}

		if err := tx.Model(&model.Booking{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"status":           model.BookingStatusCancelled,
				"cancelled_reason": "user cancelled",
			}).Error; err != nil {
			return err
		}

		for _, item := range booking.Items {
			if err := catRepo.ReleaseStock(ctx, item.TicketCategoryID, item.Quantity); err != nil {
				return err
			}
		}

		booking.Status = model.BookingStatusCancelled
		b = &booking
		return nil
	})
	return b, err
}

// ListExpired returns bookings that have expired (pending_payment past expires_at).
func (r *BookingRepository) ListExpired(ctx context.Context) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("status = ? AND expires_at < ?", model.BookingStatusPending, time.Now()).
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

// MarkExpired sets a booking status to expired and releases stock.
func (r *BookingRepository) MarkExpired(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("id = ?", id).
		Update("status", model.BookingStatusExpired).Error
}

// ListAll returns paginated bookings with optional status/date filters for admin.
func (r *BookingRepository) ListAll(ctx context.Context, status, search string, startDate, endDate *time.Time, page, limit int) ([]model.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	q := r.db.WithContext(ctx).Model(&model.Booking{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if startDate != nil {
		q = q.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		q = q.Where("created_at <= ?", endDate)
	}
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("booking_code ILIKE ?", like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var bookings []model.Booking
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Preload("Items").Find(&bookings).Error; err != nil {
		return nil, 0, err
	}
	return bookings, total, nil
}

// DashboardStats holds aggregate admin statistics.
type DashboardStats struct {
	TotalEvents   int64   `json:"total_events"`
	TotalBookings int64   `json:"total_bookings"`
	TotalRevenue  float64 `json:"total_revenue"`
	TotalUsers    int64   `json:"total_users"`
}

// GetStats returns aggregate dashboard statistics.
func (r *BookingRepository) GetStats(ctx context.Context) (*DashboardStats, error) {
	var stats DashboardStats

	// Total events
	if err := r.db.WithContext(ctx).Model(&model.Event{}).Count(&stats.TotalEvents).Error; err != nil {
		return nil, err
	}

	// Total bookings
	if err := r.db.WithContext(ctx).Model(&model.Booking{}).Count(&stats.TotalBookings).Error; err != nil {
		return nil, err
	}

	// Total revenue (sum of total_amount for paid bookings)
	if err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("status = ?", model.BookingStatusPaid).
		Scan(&stats.TotalRevenue).Error; err != nil {
		return nil, err
	}

	// Total users
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func generateBookingCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return "BK-" + hex.EncodeToString(b)
}
