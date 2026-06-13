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

type AdminHandler struct {
	bookings *repository.BookingRepository
	users    *repository.UserRepository
	events   *repository.EventRepository
}

func NewAdminHandler(
	bookings *repository.BookingRepository,
	users *repository.UserRepository,
	events *repository.EventRepository,
) *AdminHandler {
	return &AdminHandler{bookings: bookings, users: users, events: events}
}

type adminCreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type adminCreateUserResponse struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

func (h *AdminHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.bookings.GetStats(r.Context())
	if err != nil {
		httputil.Internal(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, stats)
}

func (h *AdminHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")
	var startDate, endDate *time.Time
	if s := r.URL.Query().Get("start_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err == nil {
			startDate = &t
		}
	}
	if e := r.URL.Query().Get("end_date"); e != "" {
		t, err := time.Parse("2006-01-02", e)
		if err == nil {
			endDate = &t
		}
	}

	bookings, total, err := h.bookings.ListAll(r.Context(), status, search, startDate, endDate, page, limit)
	if err != nil {
		httputil.Internal(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"data":  bookings,
		"total": total,
		"page":  page,
	})
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	search := r.URL.Query().Get("search")
	role := r.URL.Query().Get("role")

	users, total, err := h.users.ListFiltered(r.Context(), search, role, page, limit)
	if err != nil {
		httputil.Internal(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"data":  users,
		"total": total,
		"page":  page,
	})
}

func (h *AdminHandler) DetailBooking(w http.ResponseWriter, r *http.Request) {
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

	httputil.JSON(w, http.StatusOK, booking)
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req adminCreateUserRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		httputil.BadRequest(w, "email, password, and full_name are required", nil)
		return
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	passwordHash, err := model.HashPassword(req.Password)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	var phone *string
	if req.Phone != "" {
		phone = &req.Phone
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		Phone:        phone,
		Role:         role,
	}

	if err := h.users.Create(r.Context(), user); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			httputil.BadRequest(w, "user with this email already exists", nil)
			return
		}
		httputil.Internal(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, adminCreateUserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	})
}
