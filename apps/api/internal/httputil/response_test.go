package httputil

import "testing"

func TestValidateStruct_Success(t *testing.T) {
	type req struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}
	r := req{Email: "user@example.com", Password: "longpassword"}
	if errs := ValidateStruct(r); len(errs) > 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateStruct_Failures(t *testing.T) {
	type req struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}
	r := req{Email: "not-an-email", Password: "short"}
	errs := ValidateStruct(r)
	if len(errs) == 0 {
		t.Fatal("expected errors")
	}
	if errs["email"] == "" {
		t.Error("expected email error")
	}
	if errs["password"] == "" {
		t.Error("expected password error")
	}
}
