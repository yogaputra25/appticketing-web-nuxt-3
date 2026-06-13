package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ticketing/api/internal/httputil"
	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
)

type EventHandler struct {
	events *repository.EventRepository
}

func NewEventHandler(events *repository.EventRepository) *EventHandler {
	return &EventHandler{events: events}
}

// =====================================================
// DTOs
// =====================================================
type createCategoryItem struct {
	Name       string  `json:"name" validate:"required,max=100"`
	Price      float64 `json:"price" validate:"required,min=0"`
	TotalStock int     `json:"total_stock"`
	Stock      int     `json:"stock"`
	MaxPerUser int     `json:"max_per_user"`
}

func (c createCategoryItem) getTotalStock() int {
	if c.TotalStock > 0 {
		return c.TotalStock
	}
	return c.Stock
}

type createEventRequest struct {
	Title       string               `json:"title" validate:"required,min=3,max=255"`
	Description *string              `json:"description" validate:"omitempty"`
	Venue       string               `json:"venue" validate:"required,max=255"`
	StartDate   string               `json:"start_date" validate:"required"`
	EndDate     string               `json:"end_date" validate:"required"`
	BannerURL   *string              `json:"banner_url" validate:"omitempty,url"`
	Categories  []createCategoryItem `json:"categories" validate:"omitempty,dive"`
}

type updateEventRequest struct {
	Title       *string              `json:"title" validate:"omitempty,min=3,max=255"`
	Description *string              `json:"description"`
	Venue       *string              `json:"venue" validate:"omitempty,max=255"`
	StartDate   *string              `json:"start_date"`
	EndDate     *string              `json:"end_date"`
	BannerURL   *string              `json:"banner_url" validate:"omitempty,url"`
	Categories  []createCategoryItem `json:"categories" validate:"omitempty,dive"`
}

const dateLayout = "2006-01-02"

// =====================================================
// Admin handlers
// =====================================================
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}
	if fields := httputil.ValidateStruct(req); len(fields) > 0 {
		httputil.BadRequest(w, "validation failed", fields)
		return
	}

	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		httputil.BadRequest(w, "invalid start_date format, expected YYYY-MM-DD", nil)
		return
	}
	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		httputil.BadRequest(w, "invalid end_date format, expected YYYY-MM-DD", nil)
		return
	}
	if !endDate.After(startDate) {
		httputil.BadRequest(w, "end_date must be after start_date", nil)
		return
	}

	e := &model.Event{
		Title:       req.Title,
		Description: req.Description,
		Venue:       req.Venue,
		StartDate:   startDate,
		EndDate:     endDate,
		BannerURL:   req.BannerURL,
		Status:      model.EventStatusDraft,
	}
	if err := h.events.Create(r.Context(), e); err != nil {
		httputil.Internal(w, err)
		return
	}

	if len(req.Categories) > 0 {
		cats := make([]model.TicketCategory, len(req.Categories))
		for i, c := range req.Categories {
			mx := c.MaxPerUser
			if mx < 1 {
				mx = 4
			}
			cats[i] = model.TicketCategory{
				EventID:        e.ID,
				Name:           c.Name,
				Price:          c.Price,
				TotalStock:     c.getTotalStock(),
				AvailableStock: c.getTotalStock(),
				MaxPerUser:     mx,
			}
		}
		if err := h.events.BulkCreateCategories(r.Context(), cats); err != nil {
			httputil.Internal(w, err)
			return
		}
		e.Categories = cats
	}

	httputil.Created(w, e)
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httputil.BadRequest(w, "invalid event id", nil)
		return
	}

	var req updateEventRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}
	if fields := httputil.ValidateStruct(req); len(fields) > 0 {
		httputil.BadRequest(w, "validation failed", fields)
		return
	}

	existing, err := h.events.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			httputil.NotFound(w, "event not found")
			return
		}
		httputil.Internal(w, err)
		return
	}

	// Apply partial updates
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.Venue != nil {
		existing.Venue = *req.Venue
	}
	if req.StartDate != nil {
		d, err := time.Parse(dateLayout, *req.StartDate)
		if err != nil {
			httputil.BadRequest(w, "invalid start_date format, expected YYYY-MM-DD", nil)
			return
		}
		existing.StartDate = d
	}
	if req.EndDate != nil {
		d, err := time.Parse(dateLayout, *req.EndDate)
		if err != nil {
			httputil.BadRequest(w, "invalid end_date format, expected YYYY-MM-DD", nil)
			return
		}
		existing.EndDate = d
	}
	if req.BannerURL != nil {
		existing.BannerURL = req.BannerURL
	}
	if !existing.EndDate.After(existing.StartDate) {
		httputil.BadRequest(w, "end_date must be after start_date", nil)
		return
	}

	if err := h.events.Update(r.Context(), existing); err != nil {
		httputil.Internal(w, err)
		return
	}

	if req.Categories != nil {
		if err := h.events.DeleteCategoriesByEvent(r.Context(), existing.ID); err != nil {
			httputil.Internal(w, err)
			return
		}
		cats := make([]model.TicketCategory, len(req.Categories))
		for i, c := range req.Categories {
			mx := c.MaxPerUser
			if mx < 1 {
				mx = 4
			}
			cats[i] = model.TicketCategory{
				EventID:        existing.ID,
				Name:           c.Name,
				Price:          c.Price,
				TotalStock:     c.getTotalStock(),
				AvailableStock: c.getTotalStock(),
				MaxPerUser:     mx,
			}
		}
		if err := h.events.BulkCreateCategories(r.Context(), cats); err != nil {
			httputil.Internal(w, err)
			return
		}
		existing.Categories = cats
	}

	httputil.OK(w, existing)
}

func (h *EventHandler) Publish(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httputil.BadRequest(w, "invalid event id", nil)
		return
	}

	existing, err := h.events.FindByIDWithCategories(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			httputil.NotFound(w, "event not found")
			return
		}
		httputil.Internal(w, err)
		return
	}
	if existing.Status == model.EventStatusPublished {
		httputil.BadRequest(w, "event is already published", nil)
		return
	}
	if len(existing.Categories) == 0 {
		httputil.BadRequest(w, "Add at least one ticket category before publishing", nil)
		return
	}

	if err := h.events.SetStatus(r.Context(), id, model.EventStatusPublished); err != nil {
		httputil.Internal(w, err)
		return
	}
	existing.Status = model.EventStatusPublished
	httputil.OK(w, existing)
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httputil.BadRequest(w, "invalid event id", nil)
		return
	}

	// Check existing
	if _, err := h.events.FindByID(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			httputil.NotFound(w, "event not found")
			return
		}
		httputil.Internal(w, err)
		return
	}

	// Block delete if bookings exist
	bookings, err := h.events.CountBookingsByEvent(r.Context(), id)
	if err != nil {
		httputil.Internal(w, err)
		return
	}
	if bookings > 0 {
		httputil.BadRequest(w, "Cannot delete event with existing bookings", nil)
		return
	}

	if err := h.events.SoftDelete(r.Context(), id); err != nil {
		httputil.Internal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// =====================================================
// Public handlers
// =====================================================
func (h *EventHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")
	upcoming := r.URL.Query().Get("upcoming") == "true"

	// Default to published+upcoming for public
	if status == "" {
		status = model.EventStatusPublished
	}

	events, total, err := h.events.List(r.Context(), repository.EventFilter{
		Status:   status,
		Upcoming: upcoming,
		Search:   search,
		Page:     page,
		Limit:    limit,
	})
	if err != nil {
		httputil.Internal(w, err)
		return
	}
	httputil.OK(w, map[string]interface{}{
		"data":  events,
		"total": total,
		"page":  pageOr(page, 1),
		"limit": pageOr(limit, 20),
	})
}

func (h *EventHandler) ListAdmin(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	events, total, err := h.events.List(r.Context(), repository.EventFilter{
		Status: status,
		Search: search,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		httputil.Internal(w, err)
		return
	}
	httputil.OK(w, map[string]interface{}{
		"data":  events,
		"total": total,
		"page":  pageOr(page, 1),
		"limit": pageOr(limit, 20),
	})
}

func (h *EventHandler) DetailPublic(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httputil.BadRequest(w, "invalid event id", nil)
		return
	}
	e, err := h.events.FindByIDWithCategories(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			httputil.NotFound(w, "event not found")
			return
		}
		httputil.Internal(w, err)
		return
	}
	if e.Status != model.EventStatusPublished {
		httputil.NotFound(w, "event not found")
		return
	}
	httputil.OK(w, e)
}

func (h *EventHandler) DetailAdmin(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		httputil.BadRequest(w, "invalid event id", nil)
		return
	}
	e, err := h.events.FindByIDWithCategories(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			httputil.NotFound(w, "event not found")
			return
		}
		httputil.Internal(w, err)
		return
	}
	httputil.OK(w, e)
}

// =====================================================
// Helpers
// =====================================================
func parseIDParam(r *http.Request, name string) (uint64, error) {
	raw := chi.URLParam(r, name)
	return strconv.ParseUint(raw, 10, 64)
}

func pageOr(n, def int) int {
	if n < 1 {
		return def
	}
	return n
}
