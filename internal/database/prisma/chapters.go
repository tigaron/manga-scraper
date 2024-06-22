package prisma

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
)

type ChapterRepo struct {
	q *PrismaClient
}

func NewChapterRepo(prismaClient *PrismaClient) *ChapterRepo {
	return &ChapterRepo{
		q: prismaClient,
	}
}

func (c *ChapterModel) toChapter() internal.Chapter {
	provider := c.Provider()
	series := c.Series()

	var sourceURL string

	sourcePath := c.SourcePath
	if sourcePath != "" {
		sourceURL = provider.Scheme + provider.Host + sourcePath
	}

	var nextSlug, nextURL, prevSlug, prevURL string

	if c.NextSlug != "" {
		nextSlug = c.NextSlug
		if c.NextPath != "" {
			nextURL = provider.Scheme + provider.Host + c.NextPath
		}
	}

	if c.PrevSlug != "" {
		prevSlug = c.PrevSlug
		if c.PrevPath != "" {
			prevURL = provider.Scheme + provider.Host + c.PrevPath
		}
	}

	contentPaths := newStringSliceFromBytes(c.ContentPaths)

	return internal.Chapter{
		Provider:   provider.Slug,
		Series:     series.Slug,
		Slug:       c.Slug,
		Number:     c.Number,
		FullTitle:  c.FullTitle,
		ShortTitle: c.ShortTitle,
		SourceURL:  sourceURL,
		ChapterNav: &internal.ChapterNav{
			NextSlug: nextSlug,
			NextURL:  nextURL,
			PrevSlug: prevSlug,
			PrevURL:  prevURL,
		},
		ContentURLs: newContentURLsFromSlice(contentPaths, provider.Scheme+provider.Host),
		SourceHref:  c.SourceHref,
	}
}

func (c *ChapterModel) toChapterUpsert(provider *ProviderModel) internal.Chapter {
	var sourceURL string

	sourcePath := c.SourcePath
	if sourcePath != "" {
		sourceURL = provider.Scheme + provider.Host + sourcePath
	}

	var nextSlug, nextURL, prevSlug, prevURL string

	if c.NextSlug != "" {
		nextSlug = c.NextSlug
		if c.NextPath != "" {
			nextURL = provider.Scheme + provider.Host + c.NextPath
		}
	}

	if c.PrevSlug != "" {
		prevSlug = c.PrevSlug
		if c.PrevPath != "" {
			prevURL = provider.Scheme + provider.Host + c.PrevPath
		}
	}

	contentPaths := newStringSliceFromBytes(c.ContentPaths)

	return internal.Chapter{
		Provider:   provider.Slug,
		Series:     c.SeriesSlug,
		Slug:       c.Slug,
		Number:     c.Number,
		FullTitle:  c.FullTitle,
		ShortTitle: c.ShortTitle,
		SourceURL:  sourceURL,
		ChapterNav: &internal.ChapterNav{
			NextSlug: nextSlug,
			NextURL:  nextURL,
			PrevSlug: prevSlug,
			PrevURL:  prevURL,
		},
		ContentURLs: newContentURLsFromSlice(contentPaths, provider.Scheme+provider.Host),
		SourceHref:  c.SourceHref,
	}
}

func (c *ChapterModel) toBC() internal.ChapterBC {
	provider := c.Provider()
	series := c.Series()

	return internal.ChapterBC{
		Provider: internal.Breadcrumb{
			Slug:  provider.Slug,
			Title: provider.Name,
		},
		Series: internal.Breadcrumb{
			Slug:  series.Slug,
			Title: series.Title,
		},
		Chapter: internal.Breadcrumb{
			Slug:  c.Slug,
			Title: c.ShortTitle,
		},
	}
}

func (s *SeriesModel) toChapterList() []internal.Chapter {
	provider := s.Provider()
	chaptersList := s.Chapters()

	if len(chaptersList) == 0 {
		return nil
	}

	result := make([]internal.Chapter, 0, len(chaptersList))

	for i := range chaptersList {
		var sourceURL string

		sourcePath := chaptersList[i].SourcePath
		if sourcePath != "" {
			sourceURL = provider.Scheme + provider.Host + sourcePath
		}

		var nextSlug, nextURL, prevSlug, prevURL string

		if chaptersList[i].NextSlug != "" {
			nextSlug = chaptersList[i].NextSlug
			if chaptersList[i].NextPath != "" {
				nextURL = provider.Scheme + provider.Host + chaptersList[i].NextPath
			}
		}

		if chaptersList[i].PrevSlug != "" {
			prevSlug = chaptersList[i].PrevSlug
			if chaptersList[i].PrevPath != "" {
				prevURL = provider.Scheme + provider.Host + chaptersList[i].PrevPath
			}
		}
		contentPaths := newStringSliceFromBytes(chaptersList[i].ContentPaths)
		result = append(result, internal.Chapter{
			Provider:   provider.Slug,
			Series:     s.Slug,
			Slug:       chaptersList[i].Slug,
			Number:     chaptersList[i].Number,
			FullTitle:  chaptersList[i].FullTitle,
			ShortTitle: chaptersList[i].ShortTitle,
			SourceURL:  sourceURL,
			ChapterNav: &internal.ChapterNav{
				NextSlug: nextSlug,
				NextURL:  nextURL,
				PrevSlug: prevSlug,
				PrevURL:  prevURL,
			},
			ContentURLs: newContentURLsFromSlice(contentPaths, provider.Scheme+provider.Host),
			SourceHref:  chaptersList[i].SourceHref,
		})
	}

	return result
}

func (s *SeriesModel) toChapterListWithRel() internal.ChapterList {
	series := s.toSeries()

	result := internal.ChapterList{
		Series: series,
	}

	provider := s.Provider()
	chaptersList := s.Chapters()

	if len(chaptersList) == 0 {
		return result
	}

	chapters := make([]internal.Chapter, 0, len(chaptersList))

	for i := range chaptersList {
		chapters = append(chapters, internal.Chapter{
			Provider:   provider.Slug,
			Series:     s.Slug,
			Slug:       chaptersList[i].Slug,
			Number:     chaptersList[i].Number,
			ShortTitle: chaptersList[i].ShortTitle,
		})
	}

	result.Chapters = chapters

	return result
}

func (c *ChapterRepo) CreateInit(ctx context.Context, params internal.CreateInitChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterRepo.CreateInit").Finish()

	chapter, err := c.q.Chapter.CreateOne(
		Chapter.Slug.Set(params.Slug),
		Chapter.Number.Set(params.Number),
		Chapter.ShortTitle.Set(params.ShortTitle),
		Chapter.SourceHref.Set(params.SourceHref),
		Chapter.FullTitle.Set(""),
		Chapter.SourcePath.Set(""),
		Chapter.NextSlug.Set(""),
		Chapter.NextPath.Set(""),
		Chapter.PrevSlug.Set(""),
		Chapter.PrevPath.Set(""),
		Chapter.ContentPaths.Set([]byte("[]")),
		Chapter.Provider.Link(
			Provider.Slug.Equals(params.Provider),
		),
		Chapter.Series.Link(Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Series),
		)),
	).With(
		Chapter.Provider.Fetch(),
		Chapter.Series.Fetch(),
	).Exec(ctx)
	if err != nil {
		if _, ok := IsErrUniqueConstraint(err); ok {
			return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUniqueConstraint, "chapter already exists")
		}

		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to create chapter")
	}

	return chapter.toChapter(), nil
}

func (c *ChapterRepo) UpsertInit(ctx context.Context, params internal.CreateInitChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterRepo.UpsertInit").Finish()

	chapter, err := c.q.Chapter.UpsertOne(
		Chapter.ChapterUnique(
			Chapter.ProviderSlug.Equals(params.Provider),
			Chapter.SeriesSlug.Equals(params.Series),
			Chapter.Slug.Equals(params.Slug),
		),
	).Create(
		Chapter.Slug.Set(params.Slug),
		Chapter.Number.Set(params.Number),
		Chapter.ShortTitle.Set(params.ShortTitle),
		Chapter.SourceHref.Set(params.SourceHref),
		Chapter.FullTitle.Set(""),
		Chapter.SourcePath.Set(""),
		Chapter.NextSlug.Set(""),
		Chapter.NextPath.Set(""),
		Chapter.PrevSlug.Set(""),
		Chapter.PrevPath.Set(""),
		Chapter.ContentPaths.Set([]byte("[]")),
		Chapter.Provider.Link(
			Provider.Slug.Equals(params.Provider),
		),
		Chapter.Series.Link(Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Series),
		)),
	).Update(
		Chapter.Slug.Set(params.Slug),
		Chapter.Number.Set(params.Number),
		Chapter.ShortTitle.Set(params.ShortTitle),
		Chapter.SourceHref.Set(params.SourceHref),
	).Exec(ctx)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to create chapter")
	}

	provider, err := c.q.Provider.FindUnique(
		Provider.Slug.Equals(params.Provider),
	).Exec(ctx)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find provider")
	}

	return chapter.toChapterUpsert(provider), nil
}

func (c *ChapterRepo) Find(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterRepo.Find").Finish()

	chapter, err := c.q.Chapter.FindUnique(
		Chapter.ChapterUnique(
			Chapter.ProviderSlug.Equals(params.Provider),
			Chapter.SeriesSlug.Equals(params.Series),
			Chapter.Slug.Equals(params.Slug),
		),
	).With(
		Chapter.Provider.Fetch(),
		Chapter.Series.Fetch(),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrNotFound, "chapter not found")
		}

		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find chapter")
	}

	return chapter.toChapter(), nil
}

func (c *ChapterRepo) FindBC(ctx context.Context, params internal.FindChapterParams) (internal.ChapterBC, error) {
	defer newSentrySpan(ctx, "ChapterRepo.FindBC").Finish()

	chapter, err := c.q.Chapter.FindUnique(
		Chapter.ChapterUnique(
			Chapter.ProviderSlug.Equals(params.Provider),
			Chapter.SeriesSlug.Equals(params.Series),
			Chapter.Slug.Equals(params.Slug),
		),
	).Select(
		Chapter.Slug.Field(),
		Chapter.ShortTitle.Field(),
	).With(
		Chapter.Provider.Fetch().Select(
			Provider.Slug.Field(),
			Provider.Name.Field(),
		),
		Chapter.Series.Fetch().Select(
			Series.Slug.Field(),
			Series.Title.Field(),
		),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.ChapterBC{}, internal.WrapErrorf(err, internal.ErrNotFound, "chapter not found")
		}

		return internal.ChapterBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find chapter")
	}

	return chapter.toBC(), nil
}

func (c *ChapterRepo) FindLatest(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterRepo.FindLatest").Finish()

	chapter, err := c.q.Chapter.FindFirst(
		Chapter.And(
			Chapter.ProviderSlug.Equals(params.Provider),
			Chapter.SeriesSlug.Equals(params.Series),
			Chapter.NextSlug.Equals(""),
			Chapter.NextPath.Equals(""),
		),
	).OrderBy(
		Chapter.Number.Order(SortOrderDesc),
	).With(
		Chapter.Provider.Fetch(),
		Chapter.Series.Fetch(),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrNotFound, "chapter not found")
		}

		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find chapter")
	}

	return chapter.toChapter(), nil
}

func (c *ChapterRepo) Count(ctx context.Context, params internal.FindChapterParams) (int, error) {
	defer newSentrySpan(ctx, "ChapterRepo.Count").Finish()

	chapters, err := c.q.Chapter.FindMany(
		Chapter.And(
			Chapter.ProviderSlug.Equals(params.Provider),
			Chapter.SeriesSlug.Equals(params.Series),
		),
	).Select(
		Chapter.Number.Field(),
	).Exec(ctx)
	if err != nil {
		return 0, internal.WrapErrorf(err, internal.ErrUnknown, "failed to count chapters")
	}

	return len(chapters), nil
}

func (c *ChapterRepo) FindAll(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterRepo.FindAll").Finish()

	series, err := c.q.Series.FindUnique(
		Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Series),
		),
	).With(
		Series.Provider.Fetch(),
		Series.Chapters.Fetch().OrderBy(
			Chapter.Number.Order(newSortOrder(params.Order)),
		),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return nil, internal.WrapErrorf(err, internal.ErrNotFound, "series not found")
		}

		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find series")
	}

	result := series.toChapterList()

	if len(result) == 0 {
		return nil, internal.WrapErrorf(err, internal.ErrNotFound, "no chapters found")
	}

	return result, nil
}

func (c *ChapterRepo) FindListWithRel(ctx context.Context, params internal.FindChapterParams) (internal.ChapterList, error) {
	defer newSentrySpan(ctx, "ChapterRepo.FindListWithRel").Finish()

	series, err := c.q.Series.FindUnique(
		Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Series),
		),
	).With(
		Series.Provider.Fetch(),
		Series.Chapters.Fetch().Select(
			Chapter.Slug.Field(),
			Chapter.ShortTitle.Field(),
			Chapter.Number.Field(),
		).OrderBy(
			Chapter.Number.Order(newSortOrder(params.Order)),
		),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.ChapterList{}, internal.WrapErrorf(err, internal.ErrNotFound, "series not found")
		}

		return internal.ChapterList{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find series")
	}

	result := series.toChapterListWithRel()

	return result, nil
}

func (c *ChapterRepo) FindPaginated(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterRepo.FindPaginated").Finish()

	series, err := c.q.Series.FindUnique(
		Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Series),
		),
	).With(
		Series.Provider.Fetch(),
		Series.Chapters.Fetch().OrderBy(
			Chapter.Number.Order(newSortOrder(params.Order)),
		).Take(params.Size).Skip(params.Size*(params.Page-1)),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return nil, internal.WrapErrorf(err, internal.ErrNotFound, "series not found")
		}

		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find series")
	}

	result := series.toChapterList()

	if len(result) == 0 {
		return nil, internal.WrapErrorf(err, internal.ErrNotFound, "no chapters found")
	}

	return result, nil
}

func (c *ChapterRepo) UpdateInit(ctx context.Context, params internal.UpdateInitChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterRepo.UpdateInit").Finish()

	chapter, err := c.q.Chapter.FindUnique(
		Chapter.ChapterUnique(
			Chapter.ProviderSlug.Equals(params.Provider),
			Chapter.SeriesSlug.Equals(params.Series),
			Chapter.Slug.Equals(params.Slug),
		),
	).With(
		Chapter.Provider.Fetch(),
		Chapter.Series.Fetch(),
	).Update(
		Chapter.FullTitle.Set(params.FullTitle),
		Chapter.SourcePath.Set(params.SourcePath),
		Chapter.ContentPaths.Set(params.ContentPaths),
		Chapter.NextSlug.Set(params.NextSlug),
		Chapter.NextPath.Set(params.NextPath),
		Chapter.PrevSlug.Set(params.PrevSlug),
		Chapter.PrevPath.Set(params.PrevPath),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrNotFound, "chapter not found")
		}

		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to update chapter")
	}

	return chapter.toChapter(), nil
}

func (c *ChapterRepo) Delete(ctx context.Context, params internal.FindChapterParams) error {
	defer newSentrySpan(ctx, "ChapterRepo.Delete").Finish()

	_, err := c.q.Chapter.FindUnique(
		Chapter.ChapterUnique(
			Chapter.ProviderSlug.Equals(params.Provider),
			Chapter.SeriesSlug.Equals(params.Series),
			Chapter.Slug.Equals(params.Slug),
		),
	).Delete().Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.WrapErrorf(err, internal.ErrNotFound, "chapter not found")
		}

		return internal.WrapErrorf(err, internal.ErrUnknown, "failed to delete chapter")
	}

	return nil
}
