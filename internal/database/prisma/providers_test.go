package prisma

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/steebchen/prisma-client-go/engine/protocol"
	"github.com/stretchr/testify/require"
)

func createRandomProvider(t *testing.T) (ProviderModel, internal.Provider) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	expModel := ProviderModel{
		InnerProvider: InnerProvider{
			Slug:     fmt.Sprintf("test-slug-%d", r.Int()),
			Name:     fmt.Sprintf("Test Provider %d", r.Int()),
			Scheme:   "https://",
			Host:     fmt.Sprintf("example%d.com", r.Int()),
			ListPath: fmt.Sprintf("/list%d", r.Int()),
			IsActive: r.Intn(2) == 1,
		},
	}

	mock.Provider.Expect(
		client.Provider.CreateOne(
			Provider.Slug.Set(expModel.Slug),
			Provider.Name.Set(expModel.Name),
			Provider.Scheme.Set(expModel.Scheme),
			Provider.Host.Set(expModel.Host),
			Provider.ListPath.Set(expModel.ListPath),
			Provider.IsActive.Set(expModel.IsActive),
		),
	).Returns(expModel)

	expResult := internal.Provider{
		Slug:     expModel.Slug,
		Name:     expModel.Name,
		IsActive: expModel.IsActive,
		BaseURL:  expModel.Scheme + expModel.Host,
		ListURL:  expModel.Scheme + expModel.Host + expModel.ListPath,
	}

	createdProvider, err := providerRepo.Create(context.Background(), internal.ProviderParams{
		Slug:     expModel.Slug,
		Name:     expModel.Name,
		Scheme:   expModel.Scheme,
		Host:     expModel.Host,
		ListPath: expModel.ListPath,
		IsActive: &expModel.IsActive,
	})

	require.NoError(t, err)
	require.Equal(t, expResult, createdProvider)

	return expModel, expResult
}

func TestProviderRepo_Create(t *testing.T) {
	createRandomProvider(t)
}

func TestProviderRepo_Create_UniqueConstraint(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)
	isActive := true

	mock.Provider.Expect(
		client.Provider.CreateOne(
			Provider.Slug.Set("test-slug"),
			Provider.Name.Set("Test Provider"),
			Provider.Scheme.Set("https://"),
			Provider.Host.Set("example.com"),
			Provider.ListPath.Set("/list"),
			Provider.IsActive.Set(isActive),
		),
	).Errors(&protocol.UserFacingError{
		IsPanic: false,
		Message: "Unique constraint failed on the fields: (`Slug`)",
		Meta: protocol.Meta{
			Target: []interface{}{"Slug"},
		},
		ErrorCode: "P2002",
	})

	_, err := providerRepo.Create(context.Background(), internal.ProviderParams{
		Slug:     "test-slug",
		Name:     "Test Provider",
		Scheme:   "https://",
		Host:     "example.com",
		ListPath: "/list",
		IsActive: &isActive,
	})

	require.Error(t, err)
	_, ok := IsErrUniqueConstraint(err)
	require.True(t, ok)
}

func TestProviderRepo_Create_UnknownError(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)
	isActive := true

	mock.Provider.Expect(
		client.Provider.CreateOne(
			Provider.Slug.Set("test-slug"),
			Provider.Name.Set("Test Provider"),
			Provider.Scheme.Set("https://"),
			Provider.Host.Set("example.com"),
			Provider.ListPath.Set("/list"),
			Provider.IsActive.Set(isActive),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := providerRepo.Create(context.Background(), internal.ProviderParams{
		Slug:     "test-slug",
		Name:     "Test Provider",
		Scheme:   "https://",
		Host:     "example.com",
		ListPath: "/list",
		IsActive: &isActive,
	})

	require.Error(t, err)
	_, ok := IsErrUniqueConstraint(err)
	require.False(t, ok)
}

func TestProviderRepo_Find(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	expModel, expResult := createRandomProvider(t)

	mock.Provider.Expect(
		client.Provider.FindUnique(
			Provider.Slug.Equals(expModel.Slug),
		),
	).Returns(expModel)

	foundProvider, err := providerRepo.Find(context.Background(), expModel.Slug)

	require.NoError(t, err)
	require.Equal(t, expResult, foundProvider)
}

func TestProviderRepo_Find_NotFound(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	mock.Provider.Expect(
		client.Provider.FindUnique(
			Provider.Slug.Equals("non-existent-slug"),
		),
	).Errors(ErrNotFound)

	_, err := providerRepo.Find(context.Background(), "non-existent-slug")

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestProviderRepo_Find_UnknownError(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	mock.Provider.Expect(
		client.Provider.FindUnique(
			Provider.Slug.Equals("unknown-error-slug"),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := providerRepo.Find(context.Background(), "unknown-error-slug")

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func TestProviderRepo_FindAll(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	expModel1, expResult1 := createRandomProvider(t)
	expModel2, expResult2 := createRandomProvider(t)

	mock.Provider.Expect(
		client.Provider.FindMany().OrderBy(
			Provider.Slug.Order(SortOrderAsc),
		),
	).ReturnsMany([]ProviderModel{expModel1, expModel2})

	foundProviders, err := providerRepo.FindAll(context.Background(), internal.ASC)

	require.NoError(t, err)
	require.Len(t, foundProviders, 2)
	require.Equal(t, expResult1, foundProviders[0])
	require.Equal(t, expResult2, foundProviders[1])
}

func TestProviderRepo_FindAll_NotFound(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	mock.Provider.Expect(
		client.Provider.FindMany().OrderBy(
			Provider.Slug.Order(SortOrderAsc),
		),
	).ReturnsMany(nil)

	result, err := providerRepo.FindAll(context.Background(), internal.ASC)

	require.Error(t, err)
	require.Empty(t, result)
}

func TestProviderRepo_FindAll_UnknownError(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	mock.Provider.Expect(
		client.Provider.FindMany().OrderBy(
			Provider.Slug.Order(SortOrderAsc),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := providerRepo.FindAll(context.Background(), internal.ASC)

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func TestProviderRepo_Update(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)
	isActive := false

	expModel, expResult := createRandomProvider(t)

	mock.Provider.Expect(
		client.Provider.FindUnique(
			Provider.Slug.Equals(expModel.Slug),
		).Update(
			Provider.Name.Set("Updated Name"),
			Provider.Scheme.Set("http://"),
			Provider.Host.Set("updated.com"),
			Provider.ListPath.Set("/updated-list"),
			Provider.IsActive.Set(isActive),
		),
	).Returns(expModel)

	updatedProvider, err := providerRepo.Update(context.Background(), internal.ProviderParams{
		Slug:     expModel.Slug,
		Name:     "Updated Name",
		Scheme:   "http://",
		Host:     "updated.com",
		ListPath: "/updated-list",
		IsActive: &isActive,
	})

	require.NoError(t, err)
	require.Equal(t, expResult, updatedProvider)
}

func TestProviderRepo_Update_NotFound(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)
	isActive := false

	mock.Provider.Expect(
		client.Provider.FindUnique(
			Provider.Slug.Equals("non-existent-slug"),
		).Update(
			Provider.Name.Set("Updated Name"),
			Provider.Scheme.Set("http://"),
			Provider.Host.Set("updated.com"),
			Provider.ListPath.Set("/updated-list"),
			Provider.IsActive.Set(isActive),
		),
	).Errors(ErrNotFound)

	_, err := providerRepo.Update(context.Background(), internal.ProviderParams{
		Slug:     "non-existent-slug",
		Name:     "Updated Name",
		Scheme:   "http://",
		Host:     "updated.com",
		ListPath: "/updated-list",
		IsActive: &isActive,
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestProviderRepo_Update_UnknownError(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)
	isActive := false

	mock.Provider.Expect(
		client.Provider.FindUnique(
			Provider.Slug.Equals("unknown-error-slug"),
		).Update(
			Provider.Name.Set("Updated Name"),
			Provider.Scheme.Set("http://"),
			Provider.Host.Set("updated.com"),
			Provider.ListPath.Set("/updated-list"),
			Provider.IsActive.Set(isActive),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := providerRepo.Update(context.Background(), internal.ProviderParams{
		Slug:     "unknown-error-slug",
		Name:     "Updated Name",
		Scheme:   "http://",
		Host:     "updated.com",
		ListPath: "/updated-list",
		IsActive: &isActive,
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func TestProviderRepo_Delete(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	expModel, _ := createRandomProvider(t)

	mock.Provider.Expect(
		client.Provider.FindUnique(
			Provider.Slug.Equals(expModel.Slug),
		).Delete(),
	).Returns(expModel)

	err := providerRepo.Delete(context.Background(), expModel.Slug)

	require.NoError(t, err)
}

func TestProviderRepo_Delete_NotFound(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	mock.Provider.Expect(
		client.Provider.FindUnique(
			Provider.Slug.Equals("non-existent-slug"),
		).Delete(),
	).Errors(ErrNotFound)

	err := providerRepo.Delete(context.Background(), "non-existent-slug")

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestProviderRepo_Delete_UnknownError(t *testing.T) {
	client, mock, ensure := NewMock()
	defer ensure(t)

	providerRepo := NewProviderRepo(client)

	mock.Provider.Expect(
		client.Provider.FindUnique(
			Provider.Slug.Equals("unknown-error-slug"),
		).Delete(),
	).Errors(fmt.Errorf("unknown error"))

	err := providerRepo.Delete(context.Background(), "unknown-error-slug")

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}
