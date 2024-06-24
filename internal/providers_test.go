package internal

import (
	"errors"
	"testing"
)

func TestCreateProviderParams_Validate_Valid(t *testing.T) {
	params := CreateValidProviderParams()

	err := params.Validate()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateProviderParams_Validate_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*ProviderParams)
	}{
		{"InvalidSlug", func(p *ProviderParams) { p.Slug = "" }},
		{"InvalidName", func(p *ProviderParams) { p.Name = "" }},
		{"InvalidScheme", func(p *ProviderParams) { p.Scheme = "" }},
		{"InvalidHost", func(p *ProviderParams) { p.Host = "" }},
		{"InvalidListPath", func(p *ProviderParams) { p.ListPath = "" }},
	}

	for _, tc := range tests {
		tc := tc // Create a local variable and assign the value of tc to it.
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			params := CreateValidProviderParams()
			tc.modify(params)

			err := params.Validate()
			if err == nil {
				t.Errorf("expected an error, got nil")
			}

			var iErr *Error
			if !errors.As(err, &iErr) {
				t.Errorf("expected an internal Error interface, got %T", err)
			}
		})
	}
}
