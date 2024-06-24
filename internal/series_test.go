package internal

import "testing"

func TestNewSortOrder(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want SortOrder
	}{
		{"asc", "asc", ASC},
		{"desc", "desc", DESC},
		{"default", "", ASC},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := NewSortOrder(tc.s)
			if got != tc.want {
				t.Errorf("Expected %v but got %v", tc.want, got)
			}
		})
	}
}

func TestCreateInitSeriesParams_Validate(t *testing.T) {
	cases := []struct {
		name    string
		params  CreateInitSeriesParams
		wantErr bool
	}{
		{"Valid params", CreateInitSeriesParams{"provider", "slug", "title", "sourcePath"}, false},
		{"Empty provider", CreateInitSeriesParams{"", "slug", "title", "sourcePath"}, true},
		{"Empty slug", CreateInitSeriesParams{"provider", "", "title", "sourcePath"}, true},
		{"Empty title", CreateInitSeriesParams{"provider", "slug", "", "sourcePath"}, true},
		{"Empty sourcePath", CreateInitSeriesParams{"provider", "slug", "title", ""}, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.params.Validate()
			if tc.wantErr && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Did not expect an error but got one: %v", err)
			}
		})
	}
}

func TestUpdateInitSeriesParams_Validate(t *testing.T) {
	cases := []struct {
		name    string
		params  UpdateInitSeriesParams
		wantErr bool
	}{
		{"Valid params", UpdateInitSeriesParams{"provider", "slug", "thumbnailURL", "synopsis", []byte("genres")}, false},
		{"Empty provider", UpdateInitSeriesParams{"", "slug", "thumbnailURL", "synopsis", []byte("genres")}, true},
		{"Empty slug", UpdateInitSeriesParams{"provider", "", "thumbnailURL", "synopsis", []byte("genres")}, true},
		{"Empty thumbnailURL", UpdateInitSeriesParams{"provider", "slug", "", "synopsis", []byte("genres")}, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.params.Validate()
			if tc.wantErr && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Did not expect an error but got one: %v", err)
			}
		})
	}
}

func TestUpdateLatestSeriesParams_Validate(t *testing.T) {
	cases := []struct {
		name    string
		params  UpdateLatestSeriesParams
		wantErr bool
	}{
		{"Valid params", UpdateLatestSeriesParams{"provider", "slug", 1, "latestChapter"}, false},
		{"Empty provider", UpdateLatestSeriesParams{"", "slug", 1, "latestChapter"}, true},
		{"Empty slug", UpdateLatestSeriesParams{"provider", "", 1, "latestChapter"}, true},
		{"AddChapters is 0", UpdateLatestSeriesParams{"provider", "slug", 0, "latestChapter"}, true},
		{"Empty latestChapter", UpdateLatestSeriesParams{"provider", "slug", 1, ""}, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.params.Validate()
			if tc.wantErr && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Did not expect an error but got one: %v", err)
			}
		})
	}
}
