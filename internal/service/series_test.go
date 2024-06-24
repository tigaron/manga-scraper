package service

import (
	"context"
	"fmt"
	"testing"

	"fourleaves.studio/manga-scraper/internal"
	"fourleaves.studio/manga-scraper/internal/service/mock"
	"go.uber.org/mock/gomock"
)

func TestSeriesService_CreateInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSeriesRepository(ctrl)
	mockSearch := mock.NewMockSeriesSearchRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(mockRepo, mockSearch, mockLogger)

	testCases := []struct {
		name          string
		params        internal.CreateInitSeriesParams
		modifyParams  func(params *internal.CreateInitSeriesParams)
		mockReturn    func()
		expectedError bool
	}{
		{
			name:         "successful creation",
			params:       *internal.CreateValidInitSeriesParams(),
			modifyParams: func(params *internal.CreateInitSeriesParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					CreateInit(gomock.Any(), gomock.Any()).
					Return(internal.Series{Slug: "test-series", Title: "Test Series"}, nil)
				mockSearch.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "validation failure",
			params: *internal.CreateValidInitSeriesParams(),
			modifyParams: func(params *internal.CreateInitSeriesParams) {
				params.Slug = ""
			},
			mockReturn:    func() {},
			expectedError: true,
		},
		{
			name:         "repository creation error",
			params:       *internal.CreateValidInitSeriesParams(),
			modifyParams: func(params *internal.CreateInitSeriesParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					CreateInit(gomock.Any(), gomock.Any()).
					Return(internal.Series{}, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
		{
			name:         "search index error",
			params:       *internal.CreateValidInitSeriesParams(),
			modifyParams: func(params *internal.CreateInitSeriesParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					CreateInit(gomock.Any(), gomock.Any()).
					Return(internal.Series{Slug: "test-series", Title: "Test Series"}, nil)
				mockSearch.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.modifyParams(&tc.params)
			tc.mockReturn()

			_, err := service.CreateInit(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSeriesService_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearch := mock.NewMockSeriesSearchRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(nil, mockSearch, mockLogger)

	testCases := []struct {
		name          string
		q             string
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful search",
			q:    "test query",
			mockReturn: func() {
				mockSearch.EXPECT().
					Search(gomock.Any(), gomock.Any()).
					Return([]internal.Series{}, nil)
			},
			expectedError: false,
		},
		{
			name: "search error",
			q:    "test query",
			mockReturn: func() {
				mockSearch.EXPECT().
					Search(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			_, err := service.Search(context.Background(), tc.q)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSeriesService_Index(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearch := mock.NewMockSeriesSearchRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(nil, mockSearch, mockLogger)

	testCases := []struct {
		name          string
		series        []internal.Series
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful index",
			series: []internal.Series{
				{Slug: "test-series-1", Title: "Test Series 1"},
				{Slug: "test-series-2", Title: "Test Series 2"},
			},
			mockReturn: func() {
				mockSearch.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Times(2).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "search error",
			series: []internal.Series{
				{Slug: "test-series-1", Title: "Test Series 1"},
				{Slug: "test-series-2", Title: "Test Series 2"},
			},
			mockReturn: func() {
				mockSearch.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			err := service.Index(context.Background(), tc.series)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSeriesService_Find(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSeriesRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(mockRepo, nil, mockLogger)

	testCases := []struct {
		name          string
		params        internal.FindSeriesParams
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful find",
			params: internal.FindSeriesParams{
				Slug: "test-series",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					Find(gomock.Any(), gomock.Any()).
					Return(internal.Series{Slug: "test-series", Title: "Test Series"}, nil)
			},
			expectedError: false,
		},
		{
			name: "repository find error",
			params: internal.FindSeriesParams{
				Slug: "test-series",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					Find(gomock.Any(), gomock.Any()).
					Return(internal.Series{}, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			_, err := service.Find(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSeriesService_FindBC(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSeriesRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(mockRepo, nil, mockLogger)

	testCases := []struct {
		name          string
		params        internal.FindSeriesParams
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful find",
			params: internal.FindSeriesParams{
				Provider: "test-provider",
				Slug:     "test-series",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					FindBC(gomock.Any(), gomock.Any()).
					Return(internal.SeriesBC{
						Provider: internal.Breadcrumb{
							Slug:  "test-provider",
							Title: "Test Provider",
						},
						Series: internal.Breadcrumb{
							Slug:  "test-series",
							Title: "Test Series",
						},
					}, nil)
			},
			expectedError: false,
		},
		{
			name: "repository find error",
			params: internal.FindSeriesParams{
				Provider: "test-provider",
				Slug:     "test-series",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					FindBC(gomock.Any(), gomock.Any()).
					Return(internal.SeriesBC{}, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			_, err := service.FindBC(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSeriesService_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSeriesRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(mockRepo, nil, mockLogger)

	testCases := []struct {
		name          string
		params        internal.FindSeriesParams
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful find",
			params: internal.FindSeriesParams{
				Provider: "test-provider",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					FindAll(gomock.Any(), gomock.Any()).
					Return([]internal.Series{
						{Slug: "test-series-1", Title: "Test Series 1"},
						{Slug: "test-series-2", Title: "Test Series 2"},
					}, nil)
			},
			expectedError: false,
		},
		{
			name: "repository find error",
			params: internal.FindSeriesParams{
				Provider: "test-provider",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					FindAll(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			_, err := service.FindAll(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSeriesService_FindPaginated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSeriesRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(mockRepo, nil, mockLogger)

	testCases := []struct {
		name          string
		params        internal.FindSeriesParams
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful find",
			params: internal.FindSeriesParams{
				Provider: "test-provider",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					FindPaginated(gomock.Any(), gomock.Any()).
					Return([]internal.Series{
						{Slug: "test-series-1", Title: "Test Series 1"},
						{Slug: "test-series-2", Title: "Test Series 2"},
					}, nil)
			},
			expectedError: false,
		},
		{
			name: "repository find error",
			params: internal.FindSeriesParams{
				Provider: "test-provider",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					FindPaginated(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			_, err := service.FindPaginated(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSeriesService_UpdateInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSeriesRepository(ctrl)
	mockSearch := mock.NewMockSeriesSearchRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(mockRepo, mockSearch, mockLogger)

	testCases := []struct {
		name          string
		params        internal.UpdateInitSeriesParams
		modifyParams  func(params *internal.UpdateInitSeriesParams)
		mockReturn    func()
		expectedError bool
	}{
		{
			name:         "successful update",
			params:       *internal.UpdateValidInitSeriesParams(),
			modifyParams: func(params *internal.UpdateInitSeriesParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					UpdateInit(gomock.Any(), gomock.Any()).
					Return(internal.Series{Slug: "test-series", Title: "Test Series"}, nil)
				mockSearch.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "validation failure",
			params: *internal.UpdateValidInitSeriesParams(),
			modifyParams: func(params *internal.UpdateInitSeriesParams) {
				params.Slug = ""
			},
			mockReturn:    func() {},
			expectedError: true,
		},
		{
			name:         "repository update error",
			params:       *internal.UpdateValidInitSeriesParams(),
			modifyParams: func(params *internal.UpdateInitSeriesParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					UpdateInit(gomock.Any(), gomock.Any()).
					Return(internal.Series{}, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
		{
			name:         "search index error",
			params:       *internal.UpdateValidInitSeriesParams(),
			modifyParams: func(params *internal.UpdateInitSeriesParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					UpdateInit(gomock.Any(), gomock.Any()).
					Return(internal.Series{Slug: "test-series", Title: "Test Series"}, nil)
				mockSearch.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.modifyParams(&tc.params)
			tc.mockReturn()

			_, err := service.UpdateInit(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSeriesService_UpdateLatest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSeriesRepository(ctrl)
	mockSearch := mock.NewMockSeriesSearchRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(mockRepo, mockSearch, mockLogger)

	testCases := []struct {
		name          string
		params        internal.UpdateLatestSeriesParams
		modifyParams  func(params *internal.UpdateLatestSeriesParams)
		mockReturn    func()
		expectedError bool
	}{
		{
			name:         "successful update",
			params:       *internal.UpdateValidLatestSeriesParams(),
			modifyParams: func(params *internal.UpdateLatestSeriesParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					UpdateLatest(gomock.Any(), gomock.Any()).
					Return(internal.Series{Slug: "test-series", Title: "Test Series"}, nil)
				mockSearch.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "validation failure",
			params: *internal.UpdateValidLatestSeriesParams(),
			modifyParams: func(params *internal.UpdateLatestSeriesParams) {
				params.Slug = ""
			},
			mockReturn:    func() {},
			expectedError: true,
		},
		{
			name:         "repository update error",
			params:       *internal.UpdateValidLatestSeriesParams(),
			modifyParams: func(params *internal.UpdateLatestSeriesParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					UpdateLatest(gomock.Any(), gomock.Any()).
					Return(internal.Series{}, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
		{
			name:         "search index error",
			params:       *internal.UpdateValidLatestSeriesParams(),
			modifyParams: func(params *internal.UpdateLatestSeriesParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					UpdateLatest(gomock.Any(), gomock.Any()).
					Return(internal.Series{Slug: "test-series", Title: "Test Series"}, nil)
				mockSearch.EXPECT().
					Index(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.modifyParams(&tc.params)
			tc.mockReturn()

			_, err := service.UpdateLatest(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSeriesService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockSeriesRepository(ctrl)
	mockSearch := mock.NewMockSeriesSearchRepository(ctrl)
	mockLogger := mock.NewMockLogger(ctrl)

	service := NewSeriesService(mockRepo, mockSearch, mockLogger)

	testCases := []struct {
		name          string
		params        internal.FindSeriesParams
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful delete",
			params: internal.FindSeriesParams{
				Provider: "test-provider",
				Slug:     "test-series",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil)
				mockSearch.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "repository delete error",
			params: internal.FindSeriesParams{
				Provider: "test-provider",
				Slug:     "test-series",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("test error"))
			},
			expectedError: true,
		},
		{
			name: "search delete error",
			params: internal.FindSeriesParams{
				Provider: "test-provider",
				Slug:     "test-series",
			},
			mockReturn: func() {
				mockRepo.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil)
				mockSearch.EXPECT().
					Delete(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			err := service.Delete(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
