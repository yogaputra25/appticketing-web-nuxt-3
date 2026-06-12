package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/ticketing/api/internal/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user. Returns ErrUserAlreadyExists if email is taken.
func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	// Check existing
	var existing model.User
	err := r.db.WithContext(ctx).Where("email = ?", u.Email).First(&existing).Error
	if err == nil {
		return ErrUserAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		// Race: another transaction inserted with same email
		if isUniqueViolation(err) {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// Update mutates allowed fields on a user.
func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]interface{}{
			"full_name": u.FullName,
			"phone":     u.Phone,
		}).Error
}

// List returns a paginated list of users.
func (r *UserRepository) List(ctx context.Context, page, limit int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var users []model.User
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
