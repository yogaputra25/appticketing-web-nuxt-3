package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ticketing/api/internal/auth"
	"github.com/ticketing/api/internal/httputil"
	"github.com/ticketing/api/internal/repository"
	"github.com/ticketing/api/internal/service"
)

type WarHandler struct {
	war    *service.WarQueue
	events *repository.EventRepository
	cats   *repository.TicketCategoryRepository
}

func NewWarHandler(war *service.WarQueue, events *repository.EventRepository, cats *repository.TicketCategoryRepository) *WarHandler {
	return &WarHandler{war: war, events: events, cats: cats}
}

// Join handles POST /api/war/join?event_id=X
func (h *WarHandler) Join(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	eventID, err := strconv.ParseUint(r.URL.Query().Get("event_id"), 10, 64)
	if err != nil || eventID == 0 {
		httputil.BadRequest(w, "invalid or missing event_id", nil)
		return
	}

	limited, err := h.war.CheckRateLimit(r.Context(), userID, eventID)
	if err != nil {
		httputil.Internal(w, err)
		return
	}
	if limited {
		httputil.TooManyRequests(w, "Too many join attempts. Please wait before trying again.")
		return
	}

	event, err := h.events.FindByIDWithCategories(r.Context(), eventID)
	if err != nil {
		if errors.Is(err, repository.ErrEventNotFound) {
			httputil.NotFound(w, "event not found")
			return
		}
		httputil.Internal(w, err)
		return
	}
	if event.Status != "published" {
		httputil.BadRequest(w, "event is not available", nil)
		return
	}

	allSoldOut := true
	for _, cat := range event.Categories {
		if cat.AvailableStock > 0 {
			allSoldOut = false
			break
		}
	}

	if !allSoldOut {
		httputil.OK(w, map[string]interface{}{
			"redirect_to_booking": true,
			"message":             "Tickets are available. Proceed to booking.",
		})
		return
	}

	position, token, err := h.war.JoinQueue(r.Context(), userID, eventID)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	httputil.OK(w, map[string]interface{}{
		"queued":   true,
		"position": position + 1,
		"token":    token,
	})
}

// Status handles GET /api/war/status?event_id=X
func (h *WarHandler) Status(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "authentication required")
		return
	}

	eventID, err := strconv.ParseUint(r.URL.Query().Get("event_id"), 10, 64)
	if err != nil || eventID == 0 {
		httputil.BadRequest(w, "invalid or missing event_id", nil)
		return
	}

	position, total, isReady, sessionToken, err := h.war.GetQueueStatus(r.Context(), userID, eventID)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	if isReady {
		httputil.OK(w, map[string]interface{}{
			"is_ready":      true,
			"session_token": sessionToken,
			"message":       "Your turn! Proceed to booking.",
		})
		return
	}

	if position == 0 && total == 0 {
		httputil.OK(w, map[string]interface{}{
			"is_ready":  false,
			"queued":    false,
			"message":   "You are not in the queue for this event.",
		})
		return
	}

	estimatedWait := estimateWait(position)

	httputil.OK(w, map[string]interface{}{
		"is_ready":       false,
		"position":       position + 1,
		"total_in_queue": total,
		"estimated_wait": estimatedWait,
	})
}

func estimateWait(position int64) string {
	if position < 5 {
		return "less than a minute"
	}
	seconds := position * 10
	if seconds < 60 {
		return strconv.FormatInt(seconds, 10) + " seconds"
	}
	return strconv.FormatInt(seconds/60, 10) + " minutes"
}
