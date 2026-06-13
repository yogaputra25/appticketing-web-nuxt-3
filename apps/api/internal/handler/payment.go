package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ticketing/api/internal/auth"
	"github.com/ticketing/api/internal/httputil"
	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
)

type PaymentHandler struct {
	payments *repository.PaymentRepository
	bookings *repository.BookingRepository
	catRepo  *repository.TicketCategoryRepository
}

func NewPaymentHandler(
	payments *repository.PaymentRepository,
	bookings *repository.BookingRepository,
	catRepo *repository.TicketCategoryRepository,
) *PaymentHandler {
	return &PaymentHandler{payments: payments, bookings: bookings, catRepo: catRepo}
}

type createPaymentRequest struct {
	BookingID uint64 `json:"booking_id" validate:"required"`
}

type simulatePaymentRequest struct {
	Action string `json:"action" validate:"required,oneof=success fail"`
}

// Create initiates a payment for a booking.
func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	var req createPaymentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}

	booking, err := h.bookings.GetByID(r.Context(), req.BookingID)
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

	if booking.Status != model.BookingStatusPending {
		httputil.BadRequest(w, "booking is not pending payment", nil)
		return
	}

	existingPayment, err := h.payments.GetByBookingID(r.Context(), booking.ID)
	if err == nil && existingPayment != nil {
		if existingPayment.Status == model.PaymentStatusPending {
			httputil.JSON(w, http.StatusOK, existingPayment)
			return
		}
	}

	now := time.Now()
	expiresAt := now.Add(30 * time.Minute)

	payment := &model.Payment{
		PaymentCode:   generatePaymentCodeHandler(),
		BookingID:     booking.ID,
		UserID:        userID,
		Amount:        booking.TotalAmount,
		Status:        model.PaymentStatusPending,
		PaymentMethod: "simulation",
		PaidAt:        nil,
		ExpiredAt:     &expiresAt,
	}

	if err := h.payments.Create(r.Context(), payment); err != nil {
		httputil.Internal(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, payment)
}

// Simulate handles payment simulation (test mode).
func (h *PaymentHandler) Simulate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httputil.BadRequest(w, "invalid payment id", nil)
		return
	}

	var req simulatePaymentRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}

	payment, err := h.payments.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrPaymentNotFound) {
			httputil.NotFound(w, "payment not found")
			return
		}
		httputil.Internal(w, err)
		return
	}

	if payment.Status != model.PaymentStatusPending {
		httputil.BadRequest(w, "payment is not pending", nil)
		return
	}

	now := time.Now()

	if req.Action == "success" {
		eTicketCodes, err := generateETicketCodes(bookingItemsQuantity(r.Context(), h.bookings, payment.BookingID))
		if err != nil {
			httputil.Internal(w, fmt.Errorf("generate e-ticket: %w", err))
			return
		}

		if err := h.payments.UpdateStatus(r.Context(), id, model.PaymentStatusSuccess, &now, map[string]string{
			"simulated_at": now.Format(time.RFC3339),
			"result":       "approved",
		}); err != nil {
			httputil.Internal(w, err)
			return
		}

		if err := h.bookings.UpdateStatus(r.Context(), payment.BookingID, model.BookingStatusPaid); err != nil {
			httputil.Internal(w, err)
			return
		}

		if err := h.bookings.UpdateETicketCodes(r.Context(), payment.BookingID, eTicketCodes); err != nil {
			httputil.Internal(w, err)
			return
		}

		httputil.JSON(w, http.StatusOK, map[string]interface{}{
			"status":         model.PaymentStatusSuccess,
			"paid_at":        now,
			"e_ticket_codes": eTicketCodes,
		})
		return
	}

	if err := h.payments.UpdateStatus(r.Context(), id, model.PaymentStatusFailed, nil, map[string]string{
		"simulated_at": now.Format(time.RFC3339),
		"result":       "declined",
	}); err != nil {
		httputil.Internal(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"status": model.PaymentStatusFailed,
	})
}

// ListMy returns the authenticated user's payments.
func (h *PaymentHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	payments, total, err := h.payments.ListByUser(r.Context(), userID, page, limit)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"data":  payments,
		"total": total,
		"page":  page,
	})
}

// ListAll returns all payments (admin only).
func (h *PaymentHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	status := r.URL.Query().Get("status")

	payments, total, err := h.payments.ListAll(r.Context(), status, page, limit)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"data":  payments,
		"total": total,
		"page":  page,
	})
}

// helpers

func generatePaymentCodeHandler() string {
	b := make([]byte, 6)
	rand.Read(b)
	return "PAY-" + hex.EncodeToString(b)
}

func generateETicketCodes(quantity int) ([]string, error) {
	codes := make([]string, quantity)
	for i := 0; i < quantity; i++ {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		codes[i] = "TCK-" + hex.EncodeToString(b)
	}
	return codes, nil
}

func bookingItemsQuantity(ctx context.Context, bookingRepo *repository.BookingRepository, bookingID uint64) int {
	booking, err := bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return 0
	}
	total := 0
	for _, item := range booking.Items {
		total += item.Quantity
	}
	return total
}
