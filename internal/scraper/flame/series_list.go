package flame

import (
	"context"
	"sync"

	"fourleaves.studio/manga-scraper/internal"
	"fourleaves.studio/manga-scraper/internal/scraper/helper"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
)

// TODO: exclude novel
func ScrapeSeriesList(ctx context.Context, browserUrl, listUrl string, logger *zap.Logger) ([]internal.SeriesListResult, error) {
	l, err := launcher.NewManaged(browserUrl)
	if err != nil {
		return nil, err
	}

	l.Leakless(true)
	l.Headless(true)

	lC, err := l.Client()
	if err != nil {
		return nil, err
	}

	browser := rod.New().Client(lC)
	err = browser.Connect()
	if err != nil {
		return nil, err
	}

	defer browser.MustClose()

	pg, err := browser.Page(proto.TargetCreateTarget{URL: listUrl})
	if err != nil {
		return nil, err
	}

	page := pg.Context(ctx)

	var results []internal.SeriesListResult

	el, err := page.Element("div.soralist")
	if err != nil {
		return nil, err
	}

	elA, err := el.Elements("a.series")
	if err != nil {
		return nil, err
	}

	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(elA))

	for _, e := range elA {
		go func(e *rod.Element) {
			defer wg.Done()
			title, err := e.Text()
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			postId, err := e.Attribute("rel")
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			sourcePath := "/?p=" + *postId

			href, err := e.Attribute("href")
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			slug := helper.GetSlug(*href)

			mu.Lock()
			results = append(results, internal.SeriesListResult{
				Title:      title,
				Slug:       slug,
				SourcePath: sourcePath,
			})
			mu.Unlock()
		}(e)
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return nil, err
	}

	logger.Debug("Scraped series list", zap.Int("count", len(results)))

	return results, nil
}
