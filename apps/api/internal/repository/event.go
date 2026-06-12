package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/ticketing/api/internal/model"
)

var ErrEventNotFound = errors.New("event not found")

type EventFilter struct {
	Status      string
	Upcoming    bool
	Search      string
	Page        int
	Limit       int
	IncludeCats bool
}

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create inserts a new event.
func (r *EventRepository) Create(ctx context.Context, e *model.Event) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *EventRepository) FindByID(ctx context.Context, id uint64) (*model.Event, error) {
	return r.findOne(ctx, func(q *gorm.DB) *gorm.DB {
		return q.Where("id = ?", id)
	})
}

func (r *EventRepository) FindByIDWithCategories(ctx context.Context, id uint64) (*model.Event, error) {
	return r.findOne(ctx, func(q *gorm.DB) *gorm.DB {
		return q.Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Order("price DESC")
		}).Where("id = ?", id)
	})
}

func (r *EventRepository) findOne(ctx context.Context, scope func(*gorm.DB) *gorm.DB) (*model.Event, error) {
	var e model.Event
	q := r.db.WithContext(ctx)
	q = scope(q)
	if err := q.First(&e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return &e, nil
}

// Update mutates allowed fields.
func (r *EventRepository) Update(ctx context.Context, e *model.Event) error {
	return r.db.WithContext(ctx).
		Model(&model.Event{}).
		Where("id = ?", e.ID).
		Updates(map[string]interface{}{
			"title":       e.Title,
			"description": e.Description,
			"venue":       e.Venue,
			"start_date":  e.StartDate,
			"end_date":    e.EndDate,
			"banner_url":  e.BannerURL,
		}).Error
}

// SetStatus changes the event status (used for publish / cancel / finish).
func (r *EventRepository) SetStatus(ctx context.Context, id uint64, status string) error {
	res := r.db.WithContext(ctx).
		Model(&model.Event{}).
		Where("id = ?", id).
		Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrEventNotFound
	}
	return nil
}

// SoftDelete marks event as deleted.
func (r *EventRepository) SoftDelete(ctx context.Context, id uint64) error {
	res := r.db.WithContext(ctx).Delete(&model.Event{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrEventNotFound
	}
	return nil
}

// CountBookingsByEvent returns number of bookings (any status) for an event.
func (r *EventRepository) CountBookingsByEvent(ctx context.Context, eventID uint64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("event_id = ?", eventID).
		Count(&n).Error
	return n, err
}

// List returns paginated, filtered events.
func (r *EventRepository) List(ctx context.Context, f EventFilter) ([]model.Event, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}
	offset := (f.Page - 1) * f.Limit

	q := r.db.WithContext(ctx).Model(&model.Event{})
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Upcoming {
		q = q.Where("start_date > ?", "now()")
	}
	if strings.TrimSpace(f.Search) != "" {
		like := "%" + strings.ToLower(f.Search) + "%"
		q = q.Where("LOWER(title) LIKE ? OR LOWER(venue) LIKE ?", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var events []model.Event
	q = q.Order("start_date ASC").Offset(offset).Limit(f.Limit)
	if f.IncludeCats {
		q = q.Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Order("price DESC")
		})
	}
	if err := q.Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}
