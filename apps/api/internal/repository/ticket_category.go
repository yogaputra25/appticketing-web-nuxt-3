package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ticketing/api/internal/model"
)

var (
	ErrCategoryNotFound     = errors.New("ticket category not found")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrStockBelowSold       = errors.New("total stock cannot be lower than already sold")
)

type TicketCategoryRepository struct {
	db *gorm.DB
}

func NewTicketCategoryRepository(db *gorm.DB) *TicketCategoryRepository {
	return &TicketCategoryRepository{db: db}
}

// Create inserts a new category. AvailableStock defaults to TotalStock.
func (r *TicketCategoryRepository) Create(ctx context.Context, c *model.TicketCategory) error {
	if c.AvailableStock == 0 {
		c.AvailableStock = c.TotalStock
	}
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *TicketCategoryRepository) FindByID(ctx context.Context, id uint64) (*model.TicketCategory, error) {
	var c model.TicketCategory
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *TicketCategoryRepository) ListByEvent(ctx context.Context, eventID uint64, includeSoldOut bool) ([]model.TicketCategory, error) {
	q := r.db.WithContext(ctx).Where("event_id = ?", eventID).Order("price DESC")
	if !includeSoldOut {
		q = q.Where("available_stock > 0")
	}
	var cats []model.TicketCategory
	if err := q.Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

// CountCategoriesByEvent returns number of categories for an event.
func (r *TicketCategoryRepository) CountCategoriesByEvent(ctx context.Context, eventID uint64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).
		Model(&model.TicketCategory{}).
		Where("event_id = ?", eventID).
		Count(&n).Error
	return n, err
}

// Update mutates allowed fields. Validates total_stock >= sold count.
func (r *TicketCategoryRepository) Update(ctx context.Context, c *model.TicketCategory) error {
	// Load current to compute sold = total - available
	cur, err := r.FindByID(ctx, c.ID)
	if err != nil {
		return err
	}
	sold := cur.TotalStock - cur.AvailableStock
	if c.TotalStock < sold {
		return ErrStockBelowSold
	}
	// Adjust available_stock so that sold count is preserved
	c.AvailableStock = c.TotalStock - sold

	return r.db.WithContext(ctx).
		Model(&model.TicketCategory{}).
		Where("id = ?", c.ID).
		Updates(map[string]interface{}{
			"name":           c.Name,
			"description":    c.Description,
			"price":          c.Price,
			"total_stock":    c.TotalStock,
			"available_stock": c.AvailableStock,
			"max_per_user":   c.MaxPerUser,
			"sale_start_at":  c.SaleStartAt,
			"sale_end_at":    c.SaleEndAt,
		}).Error
}

// =====================================================
// STOCK RESERVATION (concurrency-safe)
// =====================================================

// GetForUpdate locks a single row using SELECT ... FOR UPDATE.
// MUST be called inside a transaction. Use ReserveStock instead —
// this is exposed for advanced cases / tests.
func (r *TicketCategoryRepository) GetForUpdate(tx *gorm.DB, id uint64) (*model.TicketCategory, error) {
	var c model.TicketCategory
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}

// ReserveStock atomically decrements available_stock by qty using a row-level
// lock to prevent overselling. Returns ErrInsufficientStock if not enough stock.
// MUST be called inside a transaction; the caller is responsible for commit/rollback.
//
// Implementation strategy:
//   1. SELECT ... FOR UPDATE on the category row (serializes concurrent reservations)
//   2. Check available_stock >= qty
//   3. UPDATE available_stock = available_stock - qty
//   4. Return the updated category
func (r *TicketCategoryRepository) ReserveStock(ctx context.Context, categoryID uint64, qty int) (*model.TicketCategory, error) {
	if qty <= 0 {
		return nil, errors.New("qty must be > 0")
	}

	var result *model.TicketCategory
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cat, err := r.GetForUpdate(tx, categoryID)
		if err != nil {
			return err
		}
		if cat.AvailableStock < qty {
			return ErrInsufficientStock
		}
		// Atomic decrement guarded by the row lock
		newAvail := cat.AvailableStock - qty
		if err := tx.Model(&model.TicketCategory{}).
			Where("id = ?", categoryID).
			Update("available_stock", newAvail).Error; err != nil {
			return err
		}
		cat.AvailableStock = newAvail
		result = cat
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ReserveStockMulti reserves stock across multiple categories atomically.
// All categories decrement in one transaction, or none do.
func (r *TicketCategoryRepository) ReserveStockMulti(ctx context.Context, items map[uint64]int) (map[uint64]*model.TicketCategory, error) {
	if len(items) == 0 {
		return map[uint64]*model.TicketCategory{}, nil
	}

	result := make(map[uint64]*model.TicketCategory, len(items))
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock all rows first (consistent lock order to prevent deadlocks)
		ids := make([]uint64, 0, len(items))
		for id := range items {
			ids = append(ids, id)
		}
		// Sort to prevent deadlock
		// (use simple sort; if you have more advanced needs use sort.Slice with locking order)
		sortUint64(ids)

		for _, id := range ids {
			cat, err := r.GetForUpdate(tx, id)
			if err != nil {
				return err
			}
			qty := items[id]
			if cat.AvailableStock < qty {
				return ErrInsufficientStock
			}
			newAvail := cat.AvailableStock - qty
			if err := tx.Model(&model.TicketCategory{}).
				Where("id = ?", id).
				Update("available_stock", newAvail).Error; err != nil {
				return err
			}
			cat.AvailableStock = newAvail
			result[id] = cat
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ReleaseStock adds qty back to available_stock (used on booking cancel/expiry).
func (r *TicketCategoryRepository) ReleaseStock(ctx context.Context, categoryID uint64, qty int) error {
	if qty <= 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&model.TicketCategory{}).
		Where("id = ?", categoryID).
		UpdateColumn("available_stock", gorm.Expr("available_stock + ?", qty)).Error
}

// IsSoldOut checks if category has no stock left.
func (r *TicketCategoryRepository) IsSoldOut(ctx context.Context, categoryID uint64) (bool, error) {
	c, err := r.FindByID(ctx, categoryID)
	if err != nil {
		return false, err
	}
	return c.AvailableStock <= 0, nil
}

// =====================================================
// Helpers
// =====================================================
func sortUint64(s []uint64) {
	// Insertion sort — small slices typical
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}

// (silence unused import 'strings' for future use)
var _ = strings.TrimSpace
