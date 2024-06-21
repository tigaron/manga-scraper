package service

import (
	"context"

	"fourleaves.studio/manga-scraper/internal"
)

type ChapterRepository interface {
	CreateInit(ctx context.Context, params internal.CreateInitChapterParams) (internal.Chapter, error)
	Find(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error)
	FindBC(ctx context.Context, params internal.FindChapterParams) (internal.ChapterBC, error)
	FindLatest(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error)
	Count(ctx context.Context, params internal.FindChapterParams) (int, error)
	FindAll(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error)
	FindPaginated(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error)
	UpdateInit(ctx context.Context, params internal.UpdateInitChapterParams) (internal.Chapter, error)
	Delete(ctx context.Context, params internal.FindChapterParams) error
}

type ChapterService struct {
	repo ChapterRepository
}

func NewChapterService(repo ChapterRepository) *ChapterService {
	return &ChapterService{
		repo: repo,
	}
}

func (s *ChapterService) CreateInit(ctx context.Context, params internal.CreateInitChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterService.CreateInit").Finish()

	if err := params.Validate(); err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrInvalidInput, "params.Validate")
	}

	chapter, err := s.repo.CreateInit(ctx, params)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.CreateInit")
	}

	return chapter, nil
}

func (s *ChapterService) Find(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterService.Find").Finish()

	chapter, err := s.repo.Find(ctx, params)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Find")
	}

	return chapter, nil
}

func (s *ChapterService) FindBC(ctx context.Context, params internal.FindChapterParams) (internal.ChapterBC, error) {
	defer newSentrySpan(ctx, "ChapterService.FindBC").Finish()

	chapter, err := s.repo.FindBC(ctx, params)
	if err != nil {
		return internal.ChapterBC{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Find")
	}

	return chapter, nil
}

func (s *ChapterService) FindLatest(ctx context.Context, params internal.FindChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterService.FindLatest").Finish()

	chapter, err := s.repo.FindLatest(ctx, params)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.FindLatest")
	}

	return chapter, nil
}

func (s *ChapterService) Count(ctx context.Context, params internal.FindChapterParams) (int, error) {
	defer newSentrySpan(ctx, "ChapterService.Count").Finish()

	count, err := s.repo.Count(ctx, params)
	if err != nil {
		return 0, internal.WrapErrorf(err, internal.ErrUnknown, "repo.Count")
	}

	return count, nil
}

func (s *ChapterService) FindAll(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterService.FindAll").Finish()

	chapters, err := s.repo.FindAll(ctx, params)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "repo.FindAll")
	}

	return chapters, nil
}

func (s *ChapterService) FindPaginated(ctx context.Context, params internal.FindChapterParams) ([]internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterService.FindPaginated").Finish()

	chapters, err := s.repo.FindPaginated(ctx, params)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrUnknown, "repo.FindPaginated")
	}

	return chapters, nil
}

func (s *ChapterService) UpdateInit(ctx context.Context, params internal.UpdateInitChapterParams) (internal.Chapter, error) {
	defer newSentrySpan(ctx, "ChapterService.UpdateInit").Finish()

	if err := params.Validate(); err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrInvalidInput, "params.Validate")
	}

	chapter, err := s.repo.UpdateInit(ctx, params)
	if err != nil {
		return internal.Chapter{}, internal.WrapErrorf(err, internal.ErrUnknown, "repo.UpdateInit")
	}

	return chapter, nil
}

func (s *ChapterService) Delete(ctx context.Context, params internal.FindChapterParams) error {
	defer newSentrySpan(ctx, "ChapterService.Delete").Finish()

	err := s.repo.Delete(ctx, params)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrUnknown, "repo.Delete")
	}

	return nil
}
