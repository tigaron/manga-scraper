package internal

import "testing"

func TestCreateScrapeRequestParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  CreateScrapeRequestParams
		wantErr bool
	}{
		{"Valid input", CreateScrapeRequestParams{Type: "SERIES_LIST", Status: "PENDING", BaseURL: "http://example.com", RequestPath: "/path", Provider: "provider"}, false},
		{"Missing type", CreateScrapeRequestParams{Status: "PENDING", BaseURL: "http://example.com", RequestPath: "/path", Provider: "provider"}, true},
		{"Missing baseURL", CreateScrapeRequestParams{Type: "SERIES_LIST", Status: "PENDING", RequestPath: "/path", Provider: "provider"}, true},
		{"Missing requestPath", CreateScrapeRequestParams{Type: "SERIES_LIST", Status: "PENDING", BaseURL: "http://example.com", Provider: "provider"}, true},
		{"Missing status", CreateScrapeRequestParams{Type: "SERIES_LIST", BaseURL: "http://example.com", RequestPath: "/path", Provider: "provider"}, true},
		{"Missing provider", CreateScrapeRequestParams{Type: "SERIES_LIST", Status: "PENDING", BaseURL: "http://example.com", RequestPath: "/path"}, true},
		{"ChapterList requires series", CreateScrapeRequestParams{Type: "CHAPTER_LIST", Status: "PENDING", BaseURL: "http://example.com", RequestPath: "/path", Provider: "provider"}, true},
		{"ChapterDetail requires series", CreateScrapeRequestParams{Type: "CHAPTER_DETAIL", Status: "PENDING", BaseURL: "http://example.com", RequestPath: "/path", Provider: "provider"}, true},
		{"ChapterDetail requires chapter", CreateScrapeRequestParams{Type: "CHAPTER_DETAIL", Status: "PENDING", BaseURL: "http://example.com", RequestPath: "/path", Provider: "provider", Series: "series"}, true},
	}

	for _, tt := range tests {
		tt := tt // Create a local variable and assign the value of tc to it.
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.params.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateScrapeRequestParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  UpdateScrapeRequestParams
		wantErr bool
	}{
		{"Valid input", UpdateScrapeRequestParams{ID: "1", Status: "COMPLETED"}, false},
		{"Missing ID", UpdateScrapeRequestParams{Status: "COMPLETED"}, true},
		{"Missing status", UpdateScrapeRequestParams{ID: "1"}, true},
	}

	for _, tt := range tests {
		tt := tt // Create a local variable and assign the value of tc to it.
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.params.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
