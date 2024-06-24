package service

import (
	"context"
	"fmt"
	"testing"

	"fourleaves.studio/manga-scraper/internal"
	"fourleaves.studio/manga-scraper/internal/service/mock"
	"go.uber.org/mock/gomock"
)

func TestProviderService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockProviderRepository(ctrl)
	service := NewProviderService(mockRepo)

	testCases := []struct {
		name          string
		params        internal.ProviderParams
		modifyParams  func(params *internal.ProviderParams)
		mockReturn    func()
		expectedError bool
	}{
		{
			name:         "successful creation",
			params:       *internal.CreateValidProviderParams(),
			modifyParams: func(params *internal.ProviderParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(internal.Provider{Slug: "test-provider", Name: "Test Provider"}, nil)
			},
			expectedError: false,
		},
		{
			name:   "validation failure",
			params: *internal.CreateValidProviderParams(),
			modifyParams: func(params *internal.ProviderParams) {
				params.Slug = ""
			},
			mockReturn:    func() {},
			expectedError: true,
		},
		{
			name:         "repository creation error",
			params:       *internal.CreateValidProviderParams(),
			modifyParams: func(params *internal.ProviderParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(internal.Provider{}, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.modifyParams(&tc.params)
			tc.mockReturn()

			_, err := service.Create(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestProviderService_Find(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockProviderRepository(ctrl)
	service := NewProviderService(mockRepo)

	testCases := []struct {
		name          string
		slug          string
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful find",
			slug: "test-provider",
			mockReturn: func() {
				mockRepo.EXPECT().
					Find(gomock.Any(), "test-provider").
					Return(internal.Provider{Slug: "test-provider", Name: "Test Provider"}, nil)
			},
			expectedError: false,
		},
		{
			name: "repository find error",
			slug: "test-provider",
			mockReturn: func() {
				mockRepo.EXPECT().
					Find(gomock.Any(), gomock.Any()).
					Return(internal.Provider{}, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			_, err := service.Find(context.Background(), tc.slug)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestProviderService_FindBC(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockProviderRepository(ctrl)
	service := NewProviderService(mockRepo)

	testCases := []struct {
		name          string
		slug          string
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful find",
			slug: "test-provider",
			mockReturn: func() {
				mockRepo.EXPECT().
					FindBC(gomock.Any(), "test-provider").
					Return(internal.ProviderBC{Provider: internal.Breadcrumb{Slug: "test-provider", Title: "Test Provider"}}, nil)
			},
			expectedError: false,
		},
		{
			name: "repository find error",
			slug: "test-provider",
			mockReturn: func() {
				mockRepo.EXPECT().
					FindBC(gomock.Any(), gomock.Any()).
					Return(internal.ProviderBC{}, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			_, err := service.FindBC(context.Background(), tc.slug)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestProviderService_FindAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockProviderRepository(ctrl)
	service := NewProviderService(mockRepo)

	testCases := []struct {
		name          string
		order         internal.SortOrder
		mockReturn    func()
		expectedError bool
	}{
		{
			name:  "successful find all",
			order: internal.ASC,
			mockReturn: func() {
				mockRepo.EXPECT().
					FindAll(gomock.Any(), internal.ASC).
					Return([]internal.Provider{{Slug: "test-provider", Name: "Test Provider"}}, nil)
			},
			expectedError: false,
		},
		{
			name:  "repository find all error",
			order: internal.ASC,
			mockReturn: func() {
				mockRepo.EXPECT().
					FindAll(gomock.Any(), internal.ASC).
					Return(nil, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			_, err := service.FindAll(context.Background(), tc.order)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestProviderService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockProviderRepository(ctrl)
	service := NewProviderService(mockRepo)

	testCases := []struct {
		name          string
		params        internal.ProviderParams
		modifyParams  func(params *internal.ProviderParams)
		mockReturn    func()
		expectedError bool
	}{
		{
			name:         "successful update",
			params:       *internal.CreateValidProviderParams(),
			modifyParams: func(params *internal.ProviderParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(internal.Provider{Slug: "test-provider", Name: "Test Provider"}, nil)
			},
			expectedError: false,
		},
		{
			name:   "validation failure",
			params: *internal.CreateValidProviderParams(),
			modifyParams: func(params *internal.ProviderParams) {
				params.Slug = ""
			},
			mockReturn:    func() {},
			expectedError: true,
		},
		{
			name:         "repository update error",
			params:       *internal.CreateValidProviderParams(),
			modifyParams: func(params *internal.ProviderParams) {},
			mockReturn: func() {
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(internal.Provider{}, fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.modifyParams(&tc.params)
			tc.mockReturn()

			_, err := service.Update(context.Background(), tc.params)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestProviderService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockProviderRepository(ctrl)
	service := NewProviderService(mockRepo)

	testCases := []struct {
		name          string
		slug          string
		mockReturn    func()
		expectedError bool
	}{
		{
			name: "successful delete",
			slug: "test-provider",
			mockReturn: func() {
				mockRepo.EXPECT().
					Delete(gomock.Any(), "test-provider").
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "repository delete error",
			slug: "test-provider",
			mockReturn: func() {
				mockRepo.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("test error"))
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockReturn()

			err := service.Delete(context.Background(), tc.slug)
			if (err != nil) != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
