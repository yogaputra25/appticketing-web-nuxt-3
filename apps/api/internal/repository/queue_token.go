package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/ticketing/api/internal/model"
)

var ErrQueueTokenNotFound = errors.New("queue token not found")

type QueueTokenRepository struct {
	db *gorm.DB
}

func NewQueueTokenRepository(db *gorm.DB) *QueueTokenRepository {
	return &QueueTokenRepository{db: db}
}

func (r *QueueTokenRepository) Create(ctx context.Context, qt *model.QueueToken) error {
	return r.db.WithContext(ctx).Create(qt).Error
}

func (r *QueueTokenRepository) FindByToken(ctx context.Context, token string) (*model.QueueToken, error) {
	var qt model.QueueToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&qt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQueueTokenNotFound
		}
		return nil, err
	}
	return &qt, nil
}

func (r *QueueTokenRepository) FindActiveByUserAndEvent(ctx context.Context, userID, eventID uint64) (*model.QueueToken, error) {
	var qt model.QueueToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND event_id = ? AND status IN ?", userID, eventID, []string{model.QueueStatusWaiting, model.QueueStatusReady}).
		First(&qt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQueueTokenNotFound
		}
		return nil, err
	}
	return &qt, nil
}

func (r *QueueTokenRepository) UpdateStatus(ctx context.Context, id uint64, status string) error {
	res := r.db.WithContext(ctx).
		Model(&model.QueueToken{}).
		Where("id = ?", id).
		Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrQueueTokenNotFound
	}
	return nil
}

func (r *QueueTokenRepository) UpdateReadyAt(ctx context.Context, id uint64, readyAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.QueueToken{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":   model.QueueStatusReady,
			"ready_at": readyAt,
		}).Error
}

func (r *QueueTokenRepository) UpdateLastActive(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.QueueToken{}).
		Where("id = ?", id).
		Update("last_active_at", time.Now()).Error
}

func (r *QueueTokenRepository) ExpireStaleTokens(ctx context.Context, idleMinutes int) (int64, error) {
	cutoff := time.Now().Add(-time.Duration(idleMinutes) * time.Minute)
	res := r.db.WithContext(ctx).
		Model(&model.QueueToken{}).
		Where("status IN ? AND last_active_at < ?", []string{model.QueueStatusWaiting, model.QueueStatusReady}, cutoff).
		Update("status", model.QueueStatusExpired)
	return res.RowsAffected, res.Error
}
