package prisma

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
)

type SeriesRepo struct {
	q *PrismaClient
}

func NewSeriesRepo(prismaClient *PrismaClient) *SeriesRepo {
	return &SeriesRepo{
		q: prismaClient,
	}
}

func (s *SeriesModel) toSeries() internal.Series {
	provider := s.Provider()

	return internal.Series{
		Provider:      provider.Slug,
		Slug:          s.Slug,
		Title:         s.Title,
		SourceURL:     provider.Scheme + provider.Host + s.SourcePath,
		CoverURL:      s.ThumbnailURL,
		Synopsis:      s.Synopsis,
		Genres:        newStringSliceFromBytes(s.Genres),
		ChaptersCount: s.ChaptersCount,
		LatestChapter: s.LatestChapter,
	}
}

func (s *SeriesModel) toBC() internal.SeriesBC {
	provider := s.Provider()

	return internal.SeriesBC{
		Provider: internal.Breadcrumb{
			Slug:  provider.Slug,
			Title: provider.Name,
		},
		Series: internal.Breadcrumb{
			Slug:  s.Slug,
			Title: s.Title,
		},
	}
}

func (p *ProviderModel) toSeriesList() []internal.Series {
	seriesList := p.Series()

	if len(seriesList) == 0 {
		return nil
	}

	result := make([]internal.Series, 0, len(seriesList))

	for i := range seriesList {
		result = append(result, internal.Series{
			Provider:      p.Slug,
			Slug:          seriesList[i].Slug,
			Title:         seriesList[i].Title,
			SourceURL:     p.Scheme + p.Host + seriesList[i].SourcePath,
			CoverURL:      seriesList[i].ThumbnailURL,
			Synopsis:      seriesList[i].Synopsis,
			Genres:        newStringSliceFromBytes(seriesList[i].Genres),
			ChaptersCount: seriesList[i].ChaptersCount,
			LatestChapter: seriesList[i].LatestChapter,
		})
	}

	return result
}

func (s *SeriesRepo) CreateInit(ctx context.Context, params internal.CreateInitSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesRepo.CreateInit").Finish()

	series, err := s.q.Series.CreateOne(
		Series.Slug.Set(params.Slug),
		Series.Title.Set(params.Title),
		Series.SourcePath.Set(""),
		Series.ThumbnailURL.Set(""),
		Series.Synopsis.Set(""),
		Series.Genres.Set([]byte("[]")),
		Series.Provider.Link(Provider.Slug.Equals(params.Provider)),
	).With(
		Series.Provider.Fetch(),
	).Exec(ctx)
	if err != nil {
		if _, ok := IsErrUniqueConstraint(err); ok {
			return internal.Series{}, internal.WrapErrorf(err, internal.ErrUniqueConstraint, "series already exists")
		}

		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to create series")
	}

	return series.toSeries(), nil
}

func (s *SeriesRepo) Find(ctx context.Context, params internal.FindSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesRepo.Find").Finish()

	series, err := s.q.Series.FindUnique(
		Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Slug),
		),
	).With(
		Series.Provider.Fetch(),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.Series{}, internal.WrapErrorf(err, internal.ErrNotFound, "series not found")
		}

		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find series")
	}

	return series.toSeries(), nil
}

func (s *SeriesRepo) FindBC(ctx context.Context, params internal.FindSeriesParams) (internal.SeriesBC, error) {
	defer newSentrySpan(ctx, "SeriesRepo.FindBC").Finish()

	series, err := s.q.Series.FindUnique(
		Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Slug),
		),
	).Select(
		Series.Slug.Field(),
		Series.Title.Field(),
	).With(
		Series.Provider.Fetch().Select(
			Provider.Slug.Field(),
			Provider.Name.Field(),
		),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.SeriesBC{}, internal.WrapErrorf(err, internal.ErrNotFound, "series not found")
		}

		return internal.SeriesBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find series")
	}

	return series.toBC(), nil
}

func (s *SeriesRepo) FindAll(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesRepo.FindAll").Finish()

	provider, err := s.q.Provider.FindUnique(
		Provider.Slug.Equals(params.Provider),
	).With(
		Provider.Series.Fetch().OrderBy(
			Series.Slug.Order(newSortOrder(params.Order)),
		),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return nil, internal.WrapErrorf(err, internal.ErrNotFound, "provider not found")
		}

		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find provider")
	}

	result := provider.toSeriesList()

	if len(result) == 0 {
		return nil, internal.WrapErrorf(err, internal.ErrNotFound, "no series found")
	}

	return result, nil
}

func (s *SeriesRepo) FindPaginated(ctx context.Context, params internal.FindSeriesParams) ([]internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesRepo.FindPaginated").Finish()

	provider, err := s.q.Provider.FindUnique(
		Provider.Slug.Equals(params.Provider),
	).With(
		Provider.Series.Fetch().OrderBy(
			Series.Slug.Order(newSortOrder(params.Order)),
		).Take(params.Size).Skip(params.Size * (params.Page - 1)),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return nil, internal.WrapErrorf(err, internal.ErrNotFound, "provider not found")
		}

		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "failed to find provider")
	}

	result := provider.toSeriesList()

	if len(result) == 0 {
		return nil, internal.WrapErrorf(err, internal.ErrNotFound, "no series found")
	}

	return result, nil
}

func (s *SeriesRepo) UpdateInit(ctx context.Context, params internal.UpdateInitSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesRepo.UpdateInit").Finish()

	series, err := s.q.Series.FindUnique(
		Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Slug),
		),
	).With(
		Series.Provider.Fetch(),
	).Update(
		Series.ThumbnailURL.Set(params.ThumbnailURL),
		Series.Synopsis.Set(params.Synopsis),
		Series.Genres.Set(params.Genres),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.Series{}, internal.WrapErrorf(err, internal.ErrNotFound, "series not found")
		}

		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to update series")
	}

	return series.toSeries(), nil
}

func (s *SeriesRepo) UpdateLatest(ctx context.Context, params internal.UpdateLatestSeriesParams) (internal.Series, error) {
	defer newSentrySpan(ctx, "SeriesRepo.UpdateLatest").Finish()

	series, err := s.q.Series.FindUnique(
		Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Slug),
		),
	).With(
		Series.Provider.Fetch(),
	).Update(
		Series.ChaptersCount.Increment(params.AddChapters),
		Series.LatestChapter.Set(params.LatestChapter),
	).Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.Series{}, internal.WrapErrorf(err, internal.ErrNotFound, "series not found")
		}

		return internal.Series{}, internal.WrapErrorf(err, internal.ErrUnknown, "failed to update series")
	}

	return series.toSeries(), nil
}

func (s *SeriesRepo) Delete(ctx context.Context, params internal.FindSeriesParams) error {
	defer newSentrySpan(ctx, "SeriesRepo.Delete").Finish()

	_, err := s.q.Series.FindUnique(
		Series.SeriesUnique(
			Series.ProviderSlug.Equals(params.Provider),
			Series.Slug.Equals(params.Slug),
		),
	).Delete().Exec(ctx)
	if err != nil {
		if IsErrNotFound(err) {
			return internal.WrapErrorf(err, internal.ErrNotFound, "series not found")
		}

		return internal.WrapErrorf(err, internal.ErrUnknown, "failed to delete series")
	}

	return nil
}
