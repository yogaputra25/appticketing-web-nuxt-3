package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ticketing/api/internal/httputil"
	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
)

type TicketCategoryHandler struct {
	cats *repository.TicketCategoryRepository
}

func NewTicketCategoryHandler(cats *repository.TicketCategoryRepository) *TicketCategoryHandler {
	return &TicketCategoryHandler{cats: cats}
}

// =====================================================
// DTOs
// =====================================================
type createCategoryRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=100"`
	Description *string    `json:"description"`
	Price       float64    `json:"price" validate:"gte=0"`
	TotalStock  int        `json:"total_stock" validate:"gt=0"`
	MaxPerUser  int        `json:"max_per_user" validate:"omitempty,gte=1,lte=20"`
	SaleStartAt *time.Time `json:"sale_start_at"`
	SaleEndAt   *time.Time `json:"sale_end_at"`
}

type updateCategoryRequest struct {
	Name        *string    `json:"name" validate:"omitempty,min=2,max=100"`
	Description *string    `json:"description"`
	Price       *float64   `json:"price" validate:"omitempty,gte=0"`
	TotalStock  *int       `json:"total_stock" validate:"omitempty,gt=0"`
	MaxPerUser  *int       `json:"max_per_user" validate:"omitempty,gte=1,lte=20"`
	SaleStartAt *time.Time `json:"sale_start_at"`
	SaleEndAt   *time.Time `json:"sale_end_at"`
}

// =====================================================
// Admin handlers
// =====================================================
func (h *TicketCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseIDParam(r, "eventId")
	if err != nil {
		httputil.BadRequest(w, "invalid event id", nil)
		return
	}

	var req createCategoryRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}
	if fields := httputil.ValidateStruct(req); len(fields) > 0 {
		httputil.BadRequest(w, "validation failed", fields)
		return
	}
	if req.SaleStartAt != nil && req.SaleEndAt != nil && req.SaleEndAt.Before(*req.SaleStartAt) {
		httputil.BadRequest(w, "sale_end_at must be after sale_start_at", nil)
		return
	}

	maxPerUser := req.MaxPerUser
	if maxPerUser == 0 {
		maxPerUser = 4
	}

	c := &model.TicketCategory{
		EventID:     eventID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		TotalStock:  req.TotalStock,
		MaxPerUser:  maxPerUser,
		SaleStartAt: req.SaleStartAt,
		SaleEndAt:   req.SaleEndAt,
	}
	if err := h.cats.Create(r.Context(), c); err != nil {
		httputil.Internal(w, err)
		return
	}
	httputil.Created(w, c)
}

func (h *TicketCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httputil.BadRequest(w, "invalid category id", nil)
		return
	}

	var req updateCategoryRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}
	if fields := httputil.ValidateStruct(req); len(fields) > 0 {
		httputil.BadRequest(w, "validation failed", fields)
		return
	}

	cur, err := h.cats.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			httputil.NotFound(w, "category not found")
			return
		}
		httputil.Internal(w, err)
		return
	}

	if req.Name != nil {
		cur.Name = *req.Name
	}
	if req.Description != nil {
		cur.Description = req.Description
	}
	if req.Price != nil {
		cur.Price = *req.Price
	}
	if req.TotalStock != nil {
		cur.TotalStock = *req.TotalStock
	}
	if req.MaxPerUser != nil {
		cur.MaxPerUser = *req.MaxPerUser
	}
	if req.SaleStartAt != nil {
		cur.SaleStartAt = req.SaleStartAt
	}
	if req.SaleEndAt != nil {
		cur.SaleEndAt = req.SaleEndAt
	}

	if err := h.cats.Update(r.Context(), cur); err != nil {
		if errors.Is(err, repository.ErrStockBelowSold) {
			httputil.BadRequest(w, err.Error(), nil)
			return
		}
		httputil.Internal(w, err)
		return
	}
	updated, _ := h.cats.FindByID(r.Context(), id)
	httputil.OK(w, updated)
}

// =====================================================
// Public handlers
// =====================================================

type categoryPublicDTO struct {
	model.TicketCategory
	IsSoldOut bool `json:"is_sold_out"`
}

func (h *TicketCategoryHandler) ListByEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseIDParam(r, "eventId")
	if err != nil {
		httputil.BadRequest(w, "invalid event id", nil)
		return
	}
	cats, err := h.cats.ListByEvent(r.Context(), eventID, true)
	if err != nil {
		httputil.Internal(w, err)
		return
	}
	out := make([]categoryPublicDTO, 0, len(cats))
	for _, c := range cats {
		out = append(out, categoryPublicDTO{TicketCategory: c, IsSoldOut: c.IsSoldOut()})
	}
	httputil.OK(w, out)
}

// =====================================================
// Helper for test
// =====================================================
func parseIntQuery(r *http.Request, name string, def int) int {
	v, _ := strconv.Atoi(r.URL.Query().Get(name))
	if v < 1 {
		return def
	}
	return v
}
