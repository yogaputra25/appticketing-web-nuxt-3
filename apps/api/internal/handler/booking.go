package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/ticketing/api/internal/auth"
	"github.com/ticketing/api/internal/httputil"
	"github.com/ticketing/api/internal/repository"
	"github.com/ticketing/api/internal/service"
)

type BookingHandler struct {
	bookings *repository.BookingRepository
	cats     *repository.TicketCategoryRepository
	war      *service.WarQueue
	ttlMin   int
}

func NewBookingHandler(bookings *repository.BookingRepository, cats *repository.TicketCategoryRepository, war *service.WarQueue, ttlMin int) *BookingHandler {
	return &BookingHandler{bookings: bookings, cats: cats, war: war, ttlMin: ttlMin}
}

type reserveRequest struct {
	EventID    uint64              `json:"event_id" validate:"required"`
	Session    string              `json:"session_token" validate:"required"`
	Items      []reserveItem       `json:"items" validate:"required,min=1,dive"`
}

type reserveItem struct {
	CategoryID uint64 `json:"category_id" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,min=1,max=4"`
}

func (h *BookingHandler) Reserve(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	var req reserveRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}
	if fields := httputil.ValidateStruct(req); len(fields) > 0 {
		httputil.BadRequest(w, "validation failed", fields)
		return
	}

	uid, _, err := h.war.ValidateBookingSession(r.Context(), req.Session)
	if err != nil {
		httputil.Unauthorized(w, "Invalid or expired session token. Please rejoin the queue.")
		return
	}
	if uid != userID {
		httputil.Forbidden(w, "session token does not match user")
		return
	}

	items := make([]repository.BookingItemInput, len(req.Items))
	for i, it := range req.Items {
		cat, err := h.cats.FindByID(r.Context(), it.CategoryID)
		if err != nil {
			if errors.Is(err, repository.ErrCategoryNotFound) {
				httputil.BadRequest(w, "ticket category not found", nil)
				return
			}
			httputil.Internal(w, err)
			return
		}
		items[i] = repository.BookingItemInput{
			CategoryID: it.CategoryID,
			Quantity:   it.Quantity,
			UnitPrice:  cat.Price,
		}
	}

	booking, err := h.bookings.Create(r.Context(), repository.CreateBookingInput{
		UserID:  userID,
		EventID: req.EventID,
		Items:   items,
	}, h.cats, h.ttlMin)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientStock) {
			httputil.Conflict(w, "Not enough stock")
			return
		}
		httputil.Internal(w, err)
		return
	}

	httputil.Created(w, booking)
}

func (h *BookingHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	bookings, total, err := h.bookings.GetByUser(r.Context(), userID, page, limit)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	httputil.OK(w, map[string]interface{}{
		"data":  bookings,
		"total": total,
		"page":  pageOr(page, 1),
		"limit": pageOr(limit, 20),
	})
}

func (h *BookingHandler) Detail(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.BadRequest(w, "invalid booking id", nil)
		return
	}

	booking, err := h.bookings.GetByIDWithAssociations(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBookingNotFound) {
			httputil.NotFound(w, "booking not found")
			return
		}
		httputil.Internal(w, err)
		return
	}

	if booking.UserID != userID {
		httputil.NotFound(w, "booking not found")
		return
	}

	httputil.OK(w, booking)
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.BadRequest(w, "invalid booking id", nil)
		return
	}

	booking, err := h.bookings.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrBookingNotFound) {
			httputil.NotFound(w, "booking not found")
			return
		}
		httputil.Internal(w, err)
		return
	}

	if booking.UserID != userID {
		httputil.NotFound(w, "booking not found")
		return
	}

	cancelled, err := h.bookings.CancelBooking(r.Context(), id, h.cats)
	if err != nil {
		if errors.Is(err, repository.ErrBookingNotPending) {
			httputil.BadRequest(w, "Paid bookings cannot be self-cancelled, contact support", nil)
			return
		}
		httputil.Internal(w, err)
		return
	}

	httputil.OK(w, cancelled)
}
