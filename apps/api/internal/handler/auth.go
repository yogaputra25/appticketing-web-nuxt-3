package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ticketing/api/internal/auth"
	"github.com/ticketing/api/internal/httputil"
	"github.com/ticketing/api/internal/model"
	"github.com/ticketing/api/internal/repository"
)

type AuthHandler struct {
	users *repository.UserRepository
	jwt   *auth.JWTManager
}

func NewAuthHandler(users *repository.UserRepository, jwt *auth.JWTManager) *AuthHandler {
	return &AuthHandler{users: users, jwt: jwt}
}

// =====================================================
// Request DTOs
// =====================================================
type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Phone    string `json:"phone" validate:"omitempty,max=32"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type updateProfileRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Phone    string `json:"phone" validate:"omitempty,max=32"`
}

// =====================================================
// Handlers
// =====================================================
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}
	if fields := httputil.ValidateStruct(req); len(fields) > 0 {
		httputil.BadRequest(w, "validation failed", fields)
		return
	}

	hash, err := model.HashPassword(req.Password)
	if err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}

	phone := strings.TrimSpace(req.Phone)
	u := &model.User{
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hash,
		FullName:     strings.TrimSpace(req.FullName),
		Phone:        nilIfEmpty(phone),
		Role:         "user",
	}

	err = h.users.Create(r.Context(), u)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			httputil.Conflict(w, "Email already registered")
			return
		}
		httputil.Internal(w, err)
		return
	}

	token, err := h.jwt.GenerateToken(u.ID, u.Role)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	httputil.Created(w, map[string]interface{}{
		"user":  u,
		"token": token,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}
	if fields := httputil.ValidateStruct(req); len(fields) > 0 {
		httputil.BadRequest(w, "validation failed", fields)
		return
	}

	u, err := h.users.FindByEmail(r.Context(), strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			httputil.Unauthorized(w, "Invalid credentials")
			return
		}
		httputil.Internal(w, err)
		return
	}

	if !u.CheckPassword(req.Password) {
		httputil.Unauthorized(w, "Invalid credentials")
		return
	}

	token, err := h.jwt.GenerateToken(u.ID, u.Role)
	if err != nil {
		httputil.Internal(w, err)
		return
	}

	httputil.OK(w, map[string]interface{}{
		"user":  u,
		"token": token,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "")
		return
	}
	u, err := h.users.FindByID(r.Context(), userID)
	if err != nil {
		httputil.NotFound(w, "user not found")
		return
	}
	httputil.OK(w, u)
}

func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httputil.Unauthorized(w, "")
		return
	}

	var req updateProfileRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.BadRequest(w, err.Error(), nil)
		return
	}
	if fields := httputil.ValidateStruct(req); len(fields) > 0 {
		httputil.BadRequest(w, "validation failed", fields)
		return
	}

	u := &model.User{
		ID:       userID,
		FullName: strings.TrimSpace(req.FullName),
		Phone:    nilIfEmpty(strings.TrimSpace(req.Phone)),
	}
	if err := h.users.Update(r.Context(), u); err != nil {
		httputil.Internal(w, err)
		return
	}

	updated, _ := h.users.FindByID(r.Context(), userID)
	httputil.OK(w, updated)
}

// =====================================================
// Helpers
// =====================================================
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
