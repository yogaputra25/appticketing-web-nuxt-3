package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/ticketing/api/internal/model"
)

var (
	ErrTicketNotFound = errors.New("ticket not found")
	ErrTicketUsed     = errors.New("ticket has already been used")
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) CreateBatch(ctx context.Context, tickets []model.Ticket) error {
	if len(tickets) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&tickets).Error
}

type TicketWithBooking struct {
	model.Ticket
	BookingCode string  `json:"booking_code"`
	EventID     uint64  `json:"event_id"`
	EventTitle  string  `json:"event_title"`
	EventVenue  string  `json:"event_venue"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

func (r *TicketRepository) ListByUser(ctx context.Context, userID uint64, page, limit int) ([]TicketWithBooking, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	q := r.db.WithContext(ctx).
		Table("tickets").
		Select(`tickets.*, bookings.booking_code, events.id as event_id, events.title as event_title,
			events.venue as event_venue, events.start_date, events.end_date`).
		Joins("JOIN bookings ON bookings.id = tickets.booking_id").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where("bookings.user_id = ?", userID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var tickets []TicketWithBooking
	if err := q.Order("tickets.created_at DESC").Offset(offset).Limit(limit).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

func (r *TicketRepository) FindByID(ctx context.Context, id uint64) (*TicketWithBooking, error) {
	var t TicketWithBooking
	if err := r.db.WithContext(ctx).
		Table("tickets").
		Select(`tickets.*, bookings.booking_code, events.id as event_id, events.title as event_title,
			events.venue as event_venue, events.start_date, events.end_date`).
		Joins("JOIN bookings ON bookings.id = tickets.booking_id").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where("tickets.id = ?", id).
		First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TicketRepository) FindByCode(ctx context.Context, code string) (*TicketWithBooking, error) {
	var t TicketWithBooking
	if err := r.db.WithContext(ctx).
		Table("tickets").
		Select(`tickets.*, bookings.booking_code, events.id as event_id, events.title as event_title,
			events.venue as event_venue, events.start_date, events.end_date`).
		Joins("JOIN bookings ON bookings.id = tickets.booking_id").
		Joins("JOIN events ON events.id = bookings.event_id").
		Where("tickets.ticket_code = ?", code).
		First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TicketRepository) MarkAsUsed(ctx context.Context, id uint64) error {
	now := time.Now()
	res := r.db.WithContext(ctx).
		Model(&model.Ticket{}).
		Where("id = ? AND status = ?", id, model.TicketStatusActive).
		Updates(map[string]interface{}{
			"status":     model.TicketStatusUsed,
			"scanned_at": now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		var t model.Ticket
		if err := r.db.WithContext(ctx).First(&t, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTicketNotFound
			}
			return err
		}
		return ErrTicketUsed
	}
	return nil
}
