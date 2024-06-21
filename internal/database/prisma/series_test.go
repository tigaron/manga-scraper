package prisma

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"fourleaves.studio/manga-scraper/internal"
	"github.com/steebchen/prisma-client-go/engine/protocol"
	"github.com/stretchr/testify/require"
)

func createRandomInitSeries(t *testing.T, provider ProviderModel) (SeriesModel, internal.Series) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	expModel := SeriesModel{
		InnerSeries: InnerSeries{
			Slug:         fmt.Sprintf("slug-%d", r.Int()),
			Title:        fmt.Sprintf("title-%d", r.Int()),
			SourcePath:   "",
			ThumbnailURL: "",
			Synopsis:     "",
			Genres:       []byte("[]"),
		},
		RelationsSeries: RelationsSeries{
			Provider: &provider,
		},
	}

	mock.Series.Expect(
		seriesRepo.q.Series.CreateOne(
			Series.Slug.Set(expModel.Slug),
			Series.Title.Set(expModel.Title),
			Series.SourcePath.Set(expModel.SourcePath),
			Series.ThumbnailURL.Set(expModel.ThumbnailURL),
			Series.Synopsis.Set(expModel.Synopsis),
			Series.Genres.Set(expModel.Genres),
			Series.Provider.Link(
				Provider.Slug.Equals(provider.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		),
	).Returns(expModel)

	expResult := internal.Series{
		Provider:      provider.Slug,
		Slug:          expModel.Slug,
		Title:         expModel.Title,
		SourceURL:     provider.Scheme + provider.Host + expModel.SourcePath,
		CoverURL:      expModel.ThumbnailURL,
		Synopsis:      expModel.Synopsis,
		Genres:        newStringSliceFromBytes(expModel.Genres),
		ChaptersCount: 0,
		LatestChapter: "",
	}

	createdSeries, err := seriesRepo.CreateInit(context.Background(), internal.CreateInitSeriesParams{
		Provider:   provider.Slug,
		Slug:       expModel.Slug,
		Title:      expModel.Title,
		SourcePath: expModel.SourcePath,
	})

	require.NoError(t, err)
	require.Equal(t, expResult, createdSeries)

	return expModel, expResult
}

func TestSeriesRepo_CreateInit(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	createRandomInitSeries(t, providerModel)
}

func TestSeriesRepo_CreateInit_UniqueConstraint(t *testing.T) {
	providerModel, _ := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.CreateOne(
			Series.Slug.Set("slug"),
			Series.Title.Set("title"),
			Series.SourcePath.Set(""),
			Series.ThumbnailURL.Set(""),
			Series.Synopsis.Set(""),
			Series.Genres.Set([]byte("[]")),
			Series.Provider.Link(
				Provider.Slug.Equals(providerModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		),
	).Errors(&protocol.UserFacingError{
		IsPanic: false,
		Message: "Unique constraint failed on the fields: (`Slug`)",
		Meta: protocol.Meta{
			Target: []interface{}{"Slug"},
		},
		ErrorCode: "P2002",
	})

	_, err := seriesRepo.CreateInit(context.Background(), internal.CreateInitSeriesParams{
		Provider:   providerModel.Slug,
		Slug:       "slug",
		Title:      "title",
		SourcePath: "",
	})

	require.Error(t, err)
	_, ok := IsErrUniqueConstraint(err)
	require.True(t, ok)
}

func TestSeriesRepo_CreateInit_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.CreateOne(
			Series.Slug.Set("slug"),
			Series.Title.Set("title"),
			Series.SourcePath.Set(""),
			Series.ThumbnailURL.Set(""),
			Series.Synopsis.Set(""),
			Series.Genres.Set([]byte("[]")),
			Series.Provider.Link(
				Provider.Slug.Equals(providerModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		),
	).Errors(fmt.Errorf("unexpected error"))

	_, err := seriesRepo.CreateInit(context.Background(), internal.CreateInitSeriesParams{
		Provider:   providerModel.Slug,
		Slug:       "slug",
		Title:      "title",
		SourcePath: "",
	})

	require.Error(t, err)
	_, ok := IsErrUniqueConstraint(err)
	require.False(t, ok)
}

func TestSeriesRepo_Find(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, series := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		),
	).Returns(seriesModel)

	foundSeries, err := seriesRepo.Find(context.Background(), internal.FindSeriesParams{
		Provider: providerModel.Slug,
		Slug:     seriesModel.Slug,
	})

	require.NoError(t, err)
	require.Equal(t, series, foundSeries)
}

func TestSeriesRepo_Find_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals("non-existent-slug"),
			),
		).With(
			Series.Provider.Fetch(),
		),
	).Errors(ErrNotFound)

	_, err := seriesRepo.Find(context.Background(), internal.FindSeriesParams{
		Provider: providerModel.Slug,
		Slug:     "non-existent-slug",
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestSeriesRepo_Find_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals("unknown-error-slug"),
			),
		).With(
			Series.Provider.Fetch(),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := seriesRepo.Find(context.Background(), internal.FindSeriesParams{
		Provider: providerModel.Slug,
		Slug:     "unknown-error-slug",
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func createRandomProviderWithSeriesRel(provider ProviderModel, series []SeriesModel) ProviderModel {
	return ProviderModel{
		InnerProvider: InnerProvider{
			Slug:     provider.Slug,
			Name:     provider.Name,
			Scheme:   provider.Scheme,
			Host:     provider.Host,
			ListPath: provider.ListPath,
			IsActive: provider.IsActive,
		},
		RelationsProvider: RelationsProvider{
			Series: series,
		},
	}
}

func TestSeriesRepo_FindAll(t *testing.T) {
	providerModel, provider := createRandomProvider(t)
	expModel1, expResult1 := createRandomInitSeries(t, providerModel)
	expModel2, expResult2 := createRandomInitSeries(t, providerModel)

	providerModelWithRel := createRandomProviderWithSeriesRel(providerModel, []SeriesModel{expModel1, expModel2})

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Provider.Expect(
		seriesRepo.q.Provider.FindUnique(
			Provider.Slug.Equals(providerModel.Slug),
		).With(
			Provider.Series.Fetch().OrderBy(
				Series.Slug.Order(SortOrderAsc),
			),
		),
	).Returns(providerModelWithRel)

	result, err := seriesRepo.FindAll(context.Background(), internal.FindSeriesParams{
		Provider: provider.Slug,
		Order:    internal.ASC,
	})

	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, expResult1, result[0])
	require.Equal(t, expResult2, result[1])
}

func TestSeriesRepo_FindAll_NotFound(t *testing.T) {
	providerModel, provider := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Provider.Expect(
		seriesRepo.q.Provider.FindUnique(
			Provider.Slug.Equals(providerModel.Slug),
		).With(
			Provider.Series.Fetch().OrderBy(
				Series.Slug.Order(SortOrderAsc),
			),
		),
	).Errors(ErrNotFound)

	_, err := seriesRepo.FindAll(context.Background(), internal.FindSeriesParams{
		Provider: provider.Slug,
		Order:    internal.ASC,
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestSeriesRepo_FindAll_UnknownError(t *testing.T) {
	providerModel, provider := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Provider.Expect(
		seriesRepo.q.Provider.FindUnique(
			Provider.Slug.Equals(providerModel.Slug),
		).With(
			Provider.Series.Fetch().OrderBy(
				Series.Slug.Order(SortOrderAsc),
			),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := seriesRepo.FindAll(context.Background(), internal.FindSeriesParams{
		Provider: provider.Slug,
		Order:    internal.ASC,
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

// TODO: Implement the following tests
// func TestSeriesRepo_FindAll_NoSeriesFound(t *testing.T) {}
// func TestSeriesRepo_UpdateStatus(t *testing.T) {}

func TestSeriesRepo_FindPaginated(t *testing.T) {
	providerModel, provider := createRandomProvider(t)
	_, _ = createRandomInitSeries(t, providerModel)
	expModel2, expResult2 := createRandomInitSeries(t, providerModel)

	providerModelWithRel := createRandomProviderWithSeriesRel(providerModel, []SeriesModel{expModel2})

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Provider.Expect(
		seriesRepo.q.Provider.FindUnique(
			Provider.Slug.Equals(providerModel.Slug),
		).With(
			Provider.Series.Fetch().OrderBy(
				Series.Slug.Order(SortOrderAsc),
			).Take(1).Skip(1),
		),
	).Returns(providerModelWithRel)

	result, err := seriesRepo.FindPaginated(context.Background(), internal.FindSeriesParams{
		Provider: provider.Slug,
		Order:    internal.ASC,
		Size:     1,
		Page:     2,
	})

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, expResult2, result[0])
}

func TestSeriesRepo_FindPaginated_NotFound(t *testing.T) {
	providerModel, provider := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Provider.Expect(
		seriesRepo.q.Provider.FindUnique(
			Provider.Slug.Equals(providerModel.Slug),
		).With(
			Provider.Series.Fetch().OrderBy(
				Series.Slug.Order(SortOrderAsc),
			).Take(1).Skip(1),
		),
	).Errors(ErrNotFound)

	_, err := seriesRepo.FindPaginated(context.Background(), internal.FindSeriesParams{
		Provider: provider.Slug,
		Order:    internal.ASC,
		Size:     1,
		Page:     2,
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestSeriesRepo_FindPaginated_UnknownError(t *testing.T) {
	providerModel, provider := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Provider.Expect(
		seriesRepo.q.Provider.FindUnique(
			Provider.Slug.Equals(providerModel.Slug),
		).With(
			Provider.Series.Fetch().OrderBy(
				Series.Slug.Order(SortOrderAsc),
			).Take(1).Skip(1),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := seriesRepo.FindPaginated(context.Background(), internal.FindSeriesParams{
		Provider: provider.Slug,
		Order:    internal.ASC,
		Size:     1,
		Page:     2,
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func TestSeriesRepo_UpdateInit(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	genresBytes := []byte("[\"genre1\",\"genre2\"]")
	var genres []string
	json.Unmarshal(genresBytes, &genres)

	updatedModel := SeriesModel{
		InnerSeries: InnerSeries{
			Slug:         seriesModel.Slug,
			Title:        seriesModel.Title,
			SourcePath:   seriesModel.SourcePath,
			ThumbnailURL: "thumbnail-url",
			Synopsis:     "synopsis",
			Genres:       genresBytes,
		},
		RelationsSeries: RelationsSeries{
			Provider: &providerModel,
		},
	}

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		).Update(
			Series.ThumbnailURL.Set(updatedModel.ThumbnailURL),
			Series.Synopsis.Set(updatedModel.Synopsis),
			Series.Genres.Set(updatedModel.Genres),
		),
	).Returns(updatedModel)

	expResult := internal.Series{
		Provider:      providerModel.Slug,
		Slug:          seriesModel.Slug,
		Title:         seriesModel.Title,
		SourceURL:     providerModel.Scheme + providerModel.Host + seriesModel.SourcePath,
		CoverURL:      updatedModel.ThumbnailURL,
		Synopsis:      updatedModel.Synopsis,
		Genres:        genres,
		ChaptersCount: 0,
		LatestChapter: "",
	}

	updatedSeries, err := seriesRepo.UpdateInit(context.Background(), internal.UpdateInitSeriesParams{
		Provider:     providerModel.Slug,
		Slug:         seriesModel.Slug,
		ThumbnailURL: "thumbnail-url",
		Synopsis:     "synopsis",
		Genres:       genresBytes,
	})

	require.NoError(t, err)
	require.Equal(t, expResult, updatedSeries)
}

func TestSeriesRepo_UpdateInit_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		).Update(
			Series.ThumbnailURL.Set("thumbnail-url"),
			Series.Synopsis.Set("synopsis"),
			Series.Genres.Set([]byte("[\"genre1\",\"genre2\"]")),
		),
	).Errors(ErrNotFound)

	_, err := seriesRepo.UpdateInit(context.Background(), internal.UpdateInitSeriesParams{
		Provider:     providerModel.Slug,
		Slug:         seriesModel.Slug,
		ThumbnailURL: "thumbnail-url",
		Synopsis:     "synopsis",
		Genres:       []byte("[\"genre1\",\"genre2\"]"),
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestSeriesRepo_UpdateInit_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		).Update(
			Series.ThumbnailURL.Set("thumbnail-url"),
			Series.Synopsis.Set("synopsis"),
			Series.Genres.Set([]byte("[\"genre1\",\"genre2\"]")),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := seriesRepo.UpdateInit(context.Background(), internal.UpdateInitSeriesParams{
		Provider:     providerModel.Slug,
		Slug:         seriesModel.Slug,
		ThumbnailURL: "thumbnail-url",
		Synopsis:     "synopsis",
		Genres:       []byte("[\"genre1\",\"genre2\"]"),
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func TestSeriesRepo_UpdateLatest(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	updatedModel := SeriesModel{
		InnerSeries: InnerSeries{
			Slug:          seriesModel.Slug,
			Title:         seriesModel.Title,
			SourcePath:    seriesModel.SourcePath,
			ThumbnailURL:  seriesModel.ThumbnailURL,
			Synopsis:      seriesModel.Synopsis,
			Genres:        seriesModel.Genres,
			ChaptersCount: 1,
			LatestChapter: "latest-chapter",
		},
		RelationsSeries: RelationsSeries{
			Provider: &providerModel,
		},
	}

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		).Update(
			Series.ChaptersCount.Increment(1),
			Series.LatestChapter.Set("latest-chapter"),
		),
	).Returns(updatedModel)

	expResult := internal.Series{
		Provider:      providerModel.Slug,
		Slug:          seriesModel.Slug,
		Title:         seriesModel.Title,
		SourceURL:     providerModel.Scheme + providerModel.Host + seriesModel.SourcePath,
		CoverURL:      seriesModel.ThumbnailURL,
		Synopsis:      seriesModel.Synopsis,
		Genres:        newStringSliceFromBytes(seriesModel.Genres),
		ChaptersCount: 1,
		LatestChapter: "latest-chapter",
	}

	updatedSeries, err := seriesRepo.UpdateLatest(context.Background(), internal.UpdateLatestSeriesParams{
		Provider:      providerModel.Slug,
		Slug:          seriesModel.Slug,
		AddChapters:   1,
		LatestChapter: "latest-chapter",
	})

	require.NoError(t, err)
	require.Equal(t, expResult, updatedSeries)
}

func TestSeriesRepo_UpdateLatest_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		).Update(
			Series.ChaptersCount.Increment(1),
			Series.LatestChapter.Set("latest-chapter"),
		),
	).Errors(ErrNotFound)

	_, err := seriesRepo.UpdateLatest(context.Background(), internal.UpdateLatestSeriesParams{
		Provider:      providerModel.Slug,
		Slug:          seriesModel.Slug,
		AddChapters:   1,
		LatestChapter: "latest-chapter",
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestSeriesRepo_UpdateLatest_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
		).Update(
			Series.ChaptersCount.Increment(1),
			Series.LatestChapter.Set("latest-chapter"),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := seriesRepo.UpdateLatest(context.Background(), internal.UpdateLatestSeriesParams{
		Provider:      providerModel.Slug,
		Slug:          seriesModel.Slug,
		AddChapters:   1,
		LatestChapter: "latest-chapter",
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func TestSeriesRepo_Delete(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).Delete(),
	).Returns(seriesModel)

	err := seriesRepo.Delete(context.Background(), internal.FindSeriesParams{
		Provider: providerModel.Slug,
		Slug:     seriesModel.Slug,
	})

	require.NoError(t, err)
}

func TestSeriesRepo_Delete_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals("non-existent-slug"),
			),
		).Delete(),
	).Errors(ErrNotFound)

	err := seriesRepo.Delete(context.Background(), internal.FindSeriesParams{
		Provider: providerModel.Slug,
		Slug:     "non-existent-slug",
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestSeriesRepo_Delete_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)

	client, mock, ensure := NewMock()
	defer ensure(t)

	seriesRepo := NewSeriesRepo(client)

	mock.Series.Expect(
		seriesRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals("unknown-error-slug"),
			),
		).Delete(),
	).Errors(fmt.Errorf("unknown error"))

	err := seriesRepo.Delete(context.Background(), internal.FindSeriesParams{
		Provider: providerModel.Slug,
		Slug:     "unknown-error-slug",
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}
