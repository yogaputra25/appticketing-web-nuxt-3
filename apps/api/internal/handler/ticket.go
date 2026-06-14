package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/ticketing/api/internal/auth"
	"github.com/ticketing/api/internal/httputil"
	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
)

type TicketHandler struct {
	tickets  *repository.TicketRepository
	bookings *repository.BookingRepository
}

func NewTicketHandler(
	tickets *repository.TicketRepository,
	bookings *repository.BookingRepository,
) *TicketHandler {
	return &TicketHandler{tickets: tickets, bookings: bookings}
}

// ListMy returns the authenticated user's tickets.
func (h *TicketHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	tickets, total, err := h.tickets.ListByUser(r.Context(), userID, page, limit)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	httputil.OK(w, map[string]interface{}{
		"data":  tickets,
		"total": total,
		"page":  pageOr(page, 1),
		"limit": pageOr(limit, 20),
	})
}

// Detail returns a single ticket by ID (ownership check).
func (h *TicketHandler) Detail(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.BadRequest(w, "invalid ticket id", nil)
		return
	}

	ticket, err := h.tickets.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTicketNotFound) {
			httputil.NotFound(w, "ticket not found")
			return
		}
		httputil.Internal(w, err)
		return
	}

	booking, err := h.bookings.GetByID(r.Context(), ticket.BookingID)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	if booking.UserID != userID {
		httputil.NotFound(w, "ticket not found")
		return
	}

	httputil.OK(w, ticket)
}

// Verify checks a ticket code and marks it as used.
func (h *TicketHandler) Verify(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		httputil.BadRequest(w, "ticket code is required", nil)
		return
	}

	ticket, err := h.tickets.FindByCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, repository.ErrTicketNotFound) {
			httputil.NotFound(w, "ticket not found")
			return
		}
		httputil.Internal(w, err)
		return
	}

	if ticket.Status != model.TicketStatusActive {
		fields := map[string]string{"status": ticket.Status}
		if ticket.ScannedAt != nil {
			fields["scanned_at"] = ticket.ScannedAt.Format("2006-01-02T15:04:05Z07:00")
		}
		httputil.BadRequest(w, "ticket has already been used", fields)
		return
	}

	if err := h.tickets.MarkAsUsed(r.Context(), ticket.ID); err != nil {
		httputil.Internal(w, err)
		return
	}

	ticket.Status = model.TicketStatusUsed
	httputil.OK(w, map[string]interface{}{
		"status":  "valid",
		"ticket":  ticket,
	})
}
