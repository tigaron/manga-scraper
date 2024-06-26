package agscomics

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

func ScrapeChapterList(ctx context.Context, browserURL, seriesURL string, logger *zap.Logger) ([]internal.ChapterListResult, error) {
	l, err := launcher.NewManaged(browserURL)
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

	pg, err := browser.Page(proto.TargetCreateTarget{URL: seriesURL})
	if err != nil {
		return nil, err
	}

	page := pg.Context(ctx)

	var results []internal.ChapterListResult

	elC, err := page.Element("div.eplister")
	if err != nil {
		return nil, err
	}

	elAs, err := elC.Elements("a")
	if err != nil {
		return nil, err
	}

	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(elAs))

	for _, e := range elAs {
		go func(e *rod.Element) {
			defer wg.Done()

			href, err := e.Attribute("href")
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			slug := helper.GetSlug(*href)

			elT, err := e.Element("span.chapternum")
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			tT, err := elT.Text()
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			title := helper.GetChapterTitle(tT)

			li, err := e.Parents("li")
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			dataNum, err := li.First().Attribute("data-num")
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}

			chapterNumber := helper.GetChapterNumber(*dataNum)

			mu.Lock()
			results = append(results, internal.ChapterListResult{
				ShortTitle: title,
				Slug:       slug,
				Number:     chapterNumber,
				Href:       *href,
			})
			mu.Unlock()
		}(e)
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return nil, err
	}

	logger.Debug("Scraped chapter list", zap.Int("count", len(results)))

	return results, nil
}
