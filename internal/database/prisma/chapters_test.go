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

func createRandomInitChapter(t *testing.T, provider ProviderModel, series SeriesModel) (ChapterModel, internal.Chapter) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	expModel := ChapterModel{
		InnerChapter: InnerChapter{
			Slug:         fmt.Sprintf("slug-chapter-%d", r.Int()),
			Number:       float64(r.Int()),
			ShortTitle:   fmt.Sprintf("short-title-%d", r.Int()),
			SourceHref:   fmt.Sprintf("source-href-%d", r.Int()),
			FullTitle:    "",
			SourcePath:   "",
			NextSlug:     "",
			NextPath:     "",
			PrevSlug:     "",
			PrevPath:     "",
			ContentPaths: []byte("[]"),
		},
		RelationsChapter: RelationsChapter{
			Provider: &provider,
			Series:   &series,
		},
	}

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.CreateOne(
			Chapter.Slug.Set(expModel.Slug),
			Chapter.Number.Set(expModel.Number),
			Chapter.ShortTitle.Set(expModel.ShortTitle),
			Chapter.SourceHref.Set(expModel.SourceHref),
			Chapter.FullTitle.Set(expModel.FullTitle),
			Chapter.SourcePath.Set(expModel.SourcePath),
			Chapter.NextSlug.Set(expModel.NextSlug),
			Chapter.NextPath.Set(expModel.NextPath),
			Chapter.PrevSlug.Set(expModel.PrevSlug),
			Chapter.PrevPath.Set(expModel.PrevPath),
			Chapter.ContentPaths.Set(expModel.ContentPaths),
			Chapter.Provider.Link(
				Provider.Slug.Equals(provider.Slug),
			),
			Chapter.Series.Link(Series.SeriesUnique(
				Series.ProviderSlug.Equals(provider.Slug),
				Series.Slug.Equals(series.Slug),
			)),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		),
	).Returns(expModel)

	expResult := expModel.toChapter()

	createdChapter, err := chapterRepo.CreateInit(context.Background(), internal.CreateInitChapterParams{
		Provider:   provider.Slug,
		Series:     series.Slug,
		Slug:       expModel.Slug,
		Number:     expModel.Number,
		ShortTitle: expModel.ShortTitle,
		SourceHref: expModel.SourceHref,
	})

	require.NoError(t, err)
	require.Equal(t, expResult, createdChapter)

	return expModel, expResult
}

func TestChapterRepo_CreateInit(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	createRandomInitChapter(t, providerModel, seriesModel)
}

func TestChapterRepo_CreateInit_UniqueConstraint(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.CreateOne(
			Chapter.Slug.Set("unique-slug"),
			Chapter.Number.Set(1),
			Chapter.ShortTitle.Set("short-title"),
			Chapter.SourceHref.Set("source-href"),
			Chapter.FullTitle.Set(""),
			Chapter.SourcePath.Set(""),
			Chapter.NextSlug.Set(""),
			Chapter.NextPath.Set(""),
			Chapter.PrevSlug.Set(""),
			Chapter.PrevPath.Set(""),
			Chapter.ContentPaths.Set([]byte("[]")),
			Chapter.Provider.Link(
				Provider.Slug.Equals(providerModel.Slug),
			),
			Chapter.Series.Link(Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			)),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		),
	).Errors(&protocol.UserFacingError{
		IsPanic: false,
		Message: "Unique constraint failed on the fields: (`Slug`)",
		Meta: protocol.Meta{
			Target: []interface{}{"Slug"},
		},
		ErrorCode: "P2002",
	})

	_, err := chapterRepo.CreateInit(context.Background(), internal.CreateInitChapterParams{
		Provider:   providerModel.Slug,
		Series:     seriesModel.Slug,
		Slug:       "unique-slug",
		Number:     1,
		ShortTitle: "short-title",
		SourceHref: "source-href",
	})

	require.Error(t, err)
	_, ok := IsErrUniqueConstraint(err)
	require.True(t, ok)
}

func TestChapterRepo_CreateIni_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.CreateOne(
			Chapter.Slug.Set("unique-slug"),
			Chapter.Number.Set(1),
			Chapter.ShortTitle.Set("short-title"),
			Chapter.SourceHref.Set("source-href"),
			Chapter.FullTitle.Set(""),
			Chapter.SourcePath.Set(""),
			Chapter.NextSlug.Set(""),
			Chapter.NextPath.Set(""),
			Chapter.PrevSlug.Set(""),
			Chapter.PrevPath.Set(""),
			Chapter.ContentPaths.Set([]byte("[]")),
			Chapter.Provider.Link(
				Provider.Slug.Equals(providerModel.Slug),
			),
			Chapter.Series.Link(Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			)),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		),
	).Errors(fmt.Errorf("unexpected error"))

	_, err := chapterRepo.CreateInit(context.Background(), internal.CreateInitChapterParams{
		Provider:   providerModel.Slug,
		Series:     seriesModel.Slug,
		Slug:       "unique-slug",
		Number:     1,
		ShortTitle: "short-title",
		SourceHref: "source-href",
	})

	require.Error(t, err)
	_, ok := IsErrUniqueConstraint(err)
	require.False(t, ok)
}

func TestChapterRepo_Find(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	chapterModel, chapter := createRandomInitChapter(t, providerModel, seriesModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindUnique(
			Chapter.ChapterUnique(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.Slug.Equals(chapterModel.Slug),
			),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		),
	).Returns(chapterModel)

	foundChapter, err := chapterRepo.Find(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Slug:     chapterModel.Slug,
	})

	require.NoError(t, err)
	require.Equal(t, chapter, foundChapter)
}

func TestChapterRepo_Find_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindUnique(
			Chapter.ChapterUnique(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.Slug.Equals("not-found"),
			),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		),
	).Errors(ErrNotFound)

	_, err := chapterRepo.Find(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Slug:     "not-found",
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestChapterRepo_Find_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	chapterModel, _ := createRandomInitChapter(t, providerModel, seriesModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindUnique(
			Chapter.ChapterUnique(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.Slug.Equals(chapterModel.Slug),
			),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := chapterRepo.Find(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Slug:     chapterModel.Slug,
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func TestChapterRepo_FindLatest(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	expModel, expResult := createRandomInitChapter(t, providerModel, seriesModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindFirst(
			Chapter.And(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.NextSlug.Equals(""),
				Chapter.NextPath.Equals(""),
			),
		).OrderBy(
			Chapter.Number.Order(SortOrderDesc),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		),
	).Returns(expModel)

	foundChapter, err := chapterRepo.FindLatest(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
	})

	require.NoError(t, err)
	require.Equal(t, expResult, foundChapter)
}

func TestChapterRepo_FindLatest_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindFirst(
			Chapter.And(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.NextSlug.Equals(""),
				Chapter.NextPath.Equals(""),
			),
		).OrderBy(
			Chapter.Number.Order(SortOrderDesc),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		),
	).Errors(ErrNotFound)

	_, err := chapterRepo.FindLatest(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestChapterRepo_FindLatest_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindFirst(
			Chapter.And(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.NextSlug.Equals(""),
				Chapter.NextPath.Equals(""),
			),
		).OrderBy(
			Chapter.Number.Order(SortOrderDesc),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := chapterRepo.FindLatest(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func createRandomSeriesWithChapterRel(provider ProviderModel, series SeriesModel, chapters []ChapterModel) SeriesModel {
	return SeriesModel{
		InnerSeries: InnerSeries{
			Slug:          series.Slug,
			Title:         series.Title,
			SourcePath:    series.SourcePath,
			ThumbnailURL:  series.ThumbnailURL,
			Synopsis:      series.Synopsis,
			Genres:        series.Genres,
			Status:        series.Status,
			ChaptersCount: series.ChaptersCount,
			LatestChapter: series.LatestChapter,
		},
		RelationsSeries: RelationsSeries{
			Provider: &provider,
			Chapters: chapters,
		},
	}
}

func TestChapterRepo_Count(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	expModel1, _ := createRandomInitChapter(t, providerModel, seriesModel)
	expModel2, _ := createRandomInitChapter(t, providerModel, seriesModel)
	expModel3, _ := createRandomInitChapter(t, providerModel, seriesModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindMany(
			Chapter.And(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
			),
		).Select(
			Chapter.Number.Field(),
		),
	).ReturnsMany([]ChapterModel{expModel1, expModel2, expModel3})

	count, err := chapterRepo.Count(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
	})

	require.NoError(t, err)
	require.Equal(t, 3, count)
}

func TestChapterRepo_Count_UnkownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindMany(
			Chapter.And(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
			),
		).Select(
			Chapter.Number.Field(),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := chapterRepo.Count(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
	})

	require.Error(t, err)
}

func TestChapterRepo_FindAll(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	expModel1, expResult1 := createRandomInitChapter(t, providerModel, seriesModel)
	expModel2, expResult2 := createRandomInitChapter(t, providerModel, seriesModel)
	expModel3, expResult3 := createRandomInitChapter(t, providerModel, seriesModel)

	seriesModelWithRel := createRandomSeriesWithChapterRel(providerModel, seriesModel, []ChapterModel{expModel1, expModel2, expModel3})

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Series.Expect(
		chapterRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
			Series.Chapters.Fetch().OrderBy(
				Chapter.Number.Order(SortOrderAsc),
			),
		),
	).Returns(seriesModelWithRel)

	foundChapters, err := chapterRepo.FindAll(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Order:    internal.ASC,
	})

	require.NoError(t, err)
	require.Len(t, foundChapters, 3)
	require.Equal(t, expResult1, foundChapters[0])
	require.Equal(t, expResult2, foundChapters[1])
	require.Equal(t, expResult3, foundChapters[2])
}

func TestChapterRepo_FindAll_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Series.Expect(
		chapterRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
			Series.Chapters.Fetch().OrderBy(
				Chapter.Number.Order(SortOrderAsc),
			),
		),
	).Errors(ErrNotFound)

	_, err := chapterRepo.FindAll(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Order:    internal.ASC,
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestChapterRepo_FindAll_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Series.Expect(
		chapterRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
			Series.Chapters.Fetch().OrderBy(
				Chapter.Number.Order(SortOrderAsc),
			),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := chapterRepo.FindAll(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Order:    internal.ASC,
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

// TODO: Implement the following tests
// func TestChapterRepo_FindAll_NoSeriesFound(t *testing.T) {}

func TestChapterRepo_FindPaginated(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	expModel3, expResult3 := createRandomInitChapter(t, providerModel, seriesModel)

	seriesModelWithRel := createRandomSeriesWithChapterRel(providerModel, seriesModel, []ChapterModel{expModel3})

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Series.Expect(
		chapterRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
			Series.Chapters.Fetch().OrderBy(
				Chapter.Number.Order(SortOrderAsc),
			).Take(2).Skip(2),
		),
	).Returns(seriesModelWithRel)

	foundChapters, err := chapterRepo.FindPaginated(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Order:    internal.ASC,
		Page:     2,
		Size:     2,
	})

	require.NoError(t, err)
	require.Len(t, foundChapters, 1)
	require.Equal(t, expResult3, foundChapters[0])
}

func TestChapterRepo_FindPaginated_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Series.Expect(
		chapterRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
			Series.Chapters.Fetch().OrderBy(
				Chapter.Number.Order(SortOrderAsc),
			).Take(2).Skip(2),
		),
	).Errors(ErrNotFound)

	_, err := chapterRepo.FindPaginated(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Order:    internal.ASC,
		Page:     2,
		Size:     2,
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestChapterRepo_FindPaginated_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Series.Expect(
		chapterRepo.q.Series.FindUnique(
			Series.SeriesUnique(
				Series.ProviderSlug.Equals(providerModel.Slug),
				Series.Slug.Equals(seriesModel.Slug),
			),
		).With(
			Series.Provider.Fetch(),
			Series.Chapters.Fetch().OrderBy(
				Chapter.Number.Order(SortOrderAsc),
			).Take(2).Skip(2),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := chapterRepo.FindPaginated(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Order:    internal.ASC,
		Page:     2,
		Size:     2,
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func TestChapterRepo_UpdateInit(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	expModel, _ := createRandomInitChapter(t, providerModel, seriesModel)

	contentUrlBytes := []byte("[\"/content-url-1\", \"/content-url-2\"]")
	var contentURLs []String
	json.Unmarshal(contentUrlBytes, &contentURLs)

	updatedModel := ChapterModel{
		InnerChapter: InnerChapter{
			Slug:         expModel.Slug,
			Number:       expModel.Number,
			ShortTitle:   expModel.ShortTitle,
			SourceHref:   expModel.SourceHref,
			FullTitle:    "new-full-title",
			SourcePath:   "/new-source-path",
			NextSlug:     "new-next-slug",
			NextPath:     "/new-next-path",
			PrevSlug:     "new-prev-slug",
			PrevPath:     "/new-prev-path",
			ContentPaths: contentUrlBytes,
		},
		RelationsChapter: RelationsChapter{
			Provider: &providerModel,
			Series:   &seriesModel,
		},
	}

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindUnique(
			Chapter.ChapterUnique(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.Slug.Equals(expModel.Slug),
			),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		).Update(
			Chapter.FullTitle.Set(updatedModel.FullTitle),
			Chapter.SourcePath.Set(updatedModel.SourcePath),
			Chapter.ContentPaths.Set(updatedModel.ContentPaths),
			Chapter.NextSlug.Set(updatedModel.NextSlug),
			Chapter.NextPath.Set(updatedModel.NextPath),
			Chapter.PrevSlug.Set(updatedModel.PrevSlug),
			Chapter.PrevPath.Set(updatedModel.PrevPath),
		),
	).Returns(updatedModel)

	updatedChapter, err := chapterRepo.UpdateInit(context.Background(), internal.UpdateInitChapterParams{
		Provider:     providerModel.Slug,
		Series:       seriesModel.Slug,
		Slug:         expModel.Slug,
		FullTitle:    updatedModel.FullTitle,
		SourcePath:   updatedModel.SourcePath,
		NextSlug:     updatedModel.NextSlug,
		NextPath:     updatedModel.NextPath,
		PrevSlug:     updatedModel.PrevSlug,
		PrevPath:     updatedModel.PrevPath,
		ContentPaths: updatedModel.ContentPaths,
	})

	expResult := internal.Chapter{
		Provider:   providerModel.Slug,
		Series:     seriesModel.Slug,
		Slug:       expModel.Slug,
		Number:     expModel.Number,
		ShortTitle: expModel.ShortTitle,
		FullTitle:  updatedModel.FullTitle,
		SourceURL:  providerModel.Scheme + providerModel.Host + updatedModel.SourcePath,
		ChapterNav: &internal.ChapterNav{
			NextSlug: updatedModel.NextSlug,
			NextURL:  providerModel.Scheme + providerModel.Host + updatedModel.NextPath,
			PrevSlug: updatedModel.PrevSlug,
			PrevURL:  providerModel.Scheme + providerModel.Host + updatedModel.PrevPath,
		},
		ContentURLs: newContentURLsFromSlice(contentURLs, providerModel.Scheme+providerModel.Host),
	}

	require.NoError(t, err)
	require.Equal(t, expResult, updatedChapter)
}

func TestChapterRepo_UpdateInit_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindUnique(
			Chapter.ChapterUnique(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.Slug.Equals("not-found"),
			),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		).Update(
			Chapter.FullTitle.Set("new-full-title"),
			Chapter.SourcePath.Set("/new-source-path"),
			Chapter.ContentPaths.Set([]byte("[]")),
			Chapter.NextSlug.Set("new-next-slug"),
			Chapter.NextPath.Set("/new-next-path"),
			Chapter.PrevSlug.Set("new-prev-slug"),
			Chapter.PrevPath.Set("/new-prev-path"),
		),
	).Errors(ErrNotFound)

	_, err := chapterRepo.UpdateInit(context.Background(), internal.UpdateInitChapterParams{
		Provider:     providerModel.Slug,
		Series:       seriesModel.Slug,
		Slug:         "not-found",
		FullTitle:    "new-full-title",
		SourcePath:   "/new-source-path",
		NextSlug:     "new-next-slug",
		NextPath:     "/new-next-path",
		PrevSlug:     "new-prev-slug",
		PrevPath:     "/new-prev-path",
		ContentPaths: []byte("[]"),
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestChapterRepo_UpdateInit_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	expModel, _ := createRandomInitChapter(t, providerModel, seriesModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindUnique(
			Chapter.ChapterUnique(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.Slug.Equals(expModel.Slug),
			),
		).With(
			Chapter.Provider.Fetch(),
			Chapter.Series.Fetch(),
		).Update(
			Chapter.FullTitle.Set("new-full-title"),
			Chapter.SourcePath.Set("/new-source-path"),
			Chapter.ContentPaths.Set([]byte("[]")),
			Chapter.NextSlug.Set("new-next-slug"),
			Chapter.NextPath.Set("/new-next-path"),
			Chapter.PrevSlug.Set("new-prev-slug"),
			Chapter.PrevPath.Set("/new-prev-path"),
		),
	).Errors(fmt.Errorf("unknown error"))

	_, err := chapterRepo.UpdateInit(context.Background(), internal.UpdateInitChapterParams{
		Provider:     providerModel.Slug,
		Series:       seriesModel.Slug,
		Slug:         expModel.Slug,
		FullTitle:    "new-full-title",
		SourcePath:   "/new-source-path",
		NextSlug:     "new-next-slug",
		NextPath:     "/new-next-path",
		PrevSlug:     "new-prev-slug",
		PrevPath:     "/new-prev-path",
		ContentPaths: []byte("[]"),
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}

func TestChapterRepo_Delete(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	expModel, _ := createRandomInitChapter(t, providerModel, seriesModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindUnique(
			Chapter.ChapterUnique(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.Slug.Equals(expModel.Slug),
			),
		).Delete(),
	).Returns(expModel)

	err := chapterRepo.Delete(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Slug:     expModel.Slug,
	})

	require.NoError(t, err)
}

func TestChapterRepo_Delete_NotFound(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindUnique(
			Chapter.ChapterUnique(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.Slug.Equals("not-found"),
			),
		).Delete(),
	).Errors(ErrNotFound)

	err := chapterRepo.Delete(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Slug:     "not-found",
	})

	require.Error(t, err)
	require.True(t, IsErrNotFound(err))
}

func TestChapterRepo_Delete_UnknownError(t *testing.T) {
	providerModel, _ := createRandomProvider(t)
	seriesModel, _ := createRandomInitSeries(t, providerModel)
	expModel, _ := createRandomInitChapter(t, providerModel, seriesModel)

	client, mock, ensure := NewMock()
	defer ensure(t)

	chapterRepo := NewChapterRepo(client)

	mock.Chapter.Expect(
		chapterRepo.q.Chapter.FindUnique(
			Chapter.ChapterUnique(
				Chapter.ProviderSlug.Equals(providerModel.Slug),
				Chapter.SeriesSlug.Equals(seriesModel.Slug),
				Chapter.Slug.Equals(expModel.Slug),
			),
		).Delete(),
	).Errors(fmt.Errorf("unknown error"))

	err := chapterRepo.Delete(context.Background(), internal.FindChapterParams{
		Provider: providerModel.Slug,
		Series:   seriesModel.Slug,
		Slug:     expModel.Slug,
	})

	require.Error(t, err)
	require.False(t, IsErrNotFound(err))
}
