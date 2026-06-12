package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error   string         `json:"error"`
	Message string         `json:"message,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type ValidationErrors map[string]string

// JSON writes a JSON response.
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// OK sends 200 with payload.
func OK(w http.ResponseWriter, payload interface{}) {
	JSON(w, http.StatusOK, payload)
}

// Created sends 201 with payload.
func Created(w http.ResponseWriter, payload interface{}) {
	JSON(w, http.StatusCreated, payload)
}

// BadRequest sends 400 with message.
func BadRequest(w http.ResponseWriter, msg string, fields map[string]string) {
	JSON(w, http.StatusBadRequest, ErrorResponse{
		Error:   "bad_request",
		Message: msg,
		Fields:  fields,
	})
}

// Unauthorized sends 401.
func Unauthorized(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "unauthorized"
	}
	JSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized", Message: msg})
}

// Forbidden sends 403.
func Forbidden(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "forbidden"
	}
	JSON(w, http.StatusForbidden, ErrorResponse{Error: "forbidden", Message: msg})
}

// NotFound sends 404.
func NotFound(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "not found"
	}
	JSON(w, http.StatusNotFound, ErrorResponse{Error: "not_found", Message: msg})
}

// Conflict sends 409.
func Conflict(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusConflict, ErrorResponse{Error: "conflict", Message: msg})
}

// Internal sends 500.
func Internal(w http.ResponseWriter, err error) {
	msg := "internal server error"
	if err != nil {
		msg = err.Error()
	}
	JSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal", Message: msg})
}

// TooManyRequests sends 429.
func TooManyRequests(w http.ResponseWriter, msg string) {
	JSON(w, http.StatusTooManyRequests, ErrorResponse{Error: "rate_limited", Message: msg})
}

// DecodeJSON reads and validates the request body.
func DecodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}

// ValidateStruct uses go-playground/validator for field validation.
func ValidateStruct(v interface{}) ValidationErrors {
	validate := validator.New()
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	var verr validator.ValidationErrors
	if !errors.As(err, &verr) {
		return ValidationErrors{"_general": err.Error()}
	}
	out := ValidationErrors{}
	for _, fe := range verr {
		out[strings.ToLower(fe.Field())] = validationMessage(fe)
	}
	return out
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "gte":
		return fmt.Sprintf("must be ≥ %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be ≤ %s", fe.Param())
	case "gt":
		return fmt.Sprintf("must be > %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	default:
		return fe.Tag()
	}
}
