package flame

import (
	"context"
	"sync"

	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	"fourleaves.studio/manga-scraper/internal/scraper/helper"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func ScrapeChapterList(ctx context.Context, browserUrl, seriesUrl string) ([]v1Model.ChapterList, error) {
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

	pg, err := browser.Page(proto.TargetCreateTarget{URL: seriesUrl})
	if err != nil {
		return nil, err
	}

	page := pg.Context(ctx)

	var results []v1Model.ChapterList

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
			results = append(results, v1Model.ChapterList{
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

	return results, nil
}
