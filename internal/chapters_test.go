package internal

import (
	"errors"
	"testing"
)

func createValidCreateInitChapterParams() *CreateInitChapterParams {
	return &CreateInitChapterParams{
		Provider:   "validProvider",
		Series:     "validSeries",
		Slug:       "validSlug",
		Number:     1,
		ShortTitle: "validShortTitle",
		SourceHref: "validSourceHref",
	}
}

func TestCreateInitChapterParams_Validate_Valid(t *testing.T) {
	params := createValidCreateInitChapterParams()

	err := params.Validate()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateInitChapterParams_Validate_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*CreateInitChapterParams)
	}{
		{"InvalidProvider", func(p *CreateInitChapterParams) { p.Provider = "" }},
		{"InvalidSeries", func(p *CreateInitChapterParams) { p.Series = "" }},
		{"InvalidSlug", func(p *CreateInitChapterParams) { p.Slug = "" }},
		{"InvalidShortTitle", func(p *CreateInitChapterParams) { p.ShortTitle = "" }},
		{"InvalidSourceHref", func(p *CreateInitChapterParams) { p.SourceHref = "" }},
	}

	for _, tc := range tests {
		tc := tc // Create a local variable and assign the value of tc to it.
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			params := createValidCreateInitChapterParams()
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

func createValidUpdateInitChapterParams() *UpdateInitChapterParams {
	return &UpdateInitChapterParams{
		Provider:   "validProvider",
		Series:     "validSeries",
		Slug:       "validSlug",
		FullTitle:  "validFullTitle",
		SourcePath: "validSourcePath",
	}
}

func TestUpdateInitChapterParams_Validate_Valid(t *testing.T) {
	params := createValidUpdateInitChapterParams()

	err := params.Validate()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUpdateInitChapterParams_Validate_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*UpdateInitChapterParams)
	}{
		{"InvalidProvider", func(p *UpdateInitChapterParams) { p.Provider = "" }},
		{"InvalidSeries", func(p *UpdateInitChapterParams) { p.Series = "" }},
		{"InvalidSlug", func(p *UpdateInitChapterParams) { p.Slug = "" }},
		{"InvalidFullTitle", func(p *UpdateInitChapterParams) { p.FullTitle = "" }},
		{"InvalidSourcePath", func(p *UpdateInitChapterParams) { p.SourcePath = "" }},
	}

	for _, tc := range tests {
		tc := tc // Create a local variable and assign the value of tc to it.
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			params := createValidUpdateInitChapterParams()
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
