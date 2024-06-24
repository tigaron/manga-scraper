package internal

import (
	"errors"
	"testing"
)

func TestWrapErrorf_NewError(t *testing.T) {
	err := WrapErrorf(nil, ErrInvalidInput, "invalid input")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "invalid input" {
		t.Errorf("expected message 'invalid input', got '%s'", err.Error())
	}

	if e, ok := err.(*Error); ok {
		if e.Code() != ErrInvalidInput {
			t.Errorf("expected error code ErrInvalidInput, got %v", e.Code())
		}
	} else {
		t.Errorf("expected *Error type, got %T", err)
	}
}

func TestWrapErrorf_WrapError(t *testing.T) {
	original := errors.New("original error")

	err := WrapErrorf(original, ErrNotFound, "not found")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, original) {
		t.Error("expected original error to be wrapped")
	}

	expectedMessage := "not found: original error"
	if err.Error() != expectedMessage {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapErrorf_FormatMessage(t *testing.T) {
	expectedMessage := "error 404: not found"

	err := WrapErrorf(nil, ErrUnknown, "error %d: %s", 404, "not found")
	if err.Error() != expectedMessage {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, err.Error())
	}
}

func TestWrapErrorf_ErrorCodes(t *testing.T) {
	codes := []ErrorCode{ErrUnknown, ErrNotFound, ErrUniqueConstraint, ErrInvalidInput}
	for _, code := range codes {
		err := WrapErrorf(nil, code, "test")
		if e, ok := err.(*Error); ok {
			if e.Code() != code {
				t.Errorf("expected error code %v, got %v", code, e.Code())
			}
		} else {
			t.Errorf("expected *Error type, got %T", err)
		}
	}
}
