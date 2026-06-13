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
	ErrPaymentNotFound    = errors.New("payment not found")
	ErrPaymentNotPending  = errors.New("payment is not in pending state")
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, p *model.Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PaymentRepository) GetByID(ctx context.Context, id uint64) (*model.Payment, error) {
	var p model.Payment
	if err := r.db.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepository) GetByBookingID(ctx context.Context, bookingID uint64) (*model.Payment, error) {
	var p model.Payment
	if err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id uint64, status string, paidAt *time.Time, gatewayResp interface{}) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if paidAt != nil {
		updates["paid_at"] = paidAt
	}
	if gatewayResp != nil {
		updates["gateway_response"] = gatewayResp
	}

	res := r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("id = ?", id).
		Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrPaymentNotFound
	}
	return nil
}

func (r *PaymentRepository) ListByUser(ctx context.Context, userID uint64, page, limit int) ([]model.Payment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	q := r.db.WithContext(ctx).Model(&model.Payment{}).Where("user_id = ?", userID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var payments []model.Payment
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error; err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}

func (r *PaymentRepository) ListAll(ctx context.Context, status string, page, limit int) ([]model.Payment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	q := r.db.WithContext(ctx).Model(&model.Payment{})
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var payments []model.Payment
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error; err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}

func (r *PaymentRepository) ListExpired(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment
	if err := r.db.WithContext(ctx).
		Where("status = ? AND expired_at < ?", model.PaymentStatusPending, time.Now()).
		Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func generatePaymentCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return "PAY-" + hex.EncodeToString(b)
}
