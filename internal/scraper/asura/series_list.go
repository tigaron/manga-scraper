package asura

import (
	"context"

	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	"fourleaves.studio/manga-scraper/internal/scraper/helper"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func ScrapeSeriesList(ctx context.Context, browserUrl, listUrl string) ([]v1Model.SeriesList, error) {
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

	var results []v1Model.SeriesList

	el, err := page.Element("div.soralist")
	if err != nil {
		return nil, err
	}

	elA, err := el.Elements("a.series")
	if err != nil {
		return nil, err
	}

	var loopErr error

	for _, e := range elA {
		title, err := e.Text()
		if err != nil {
			loopErr = err
			break
		}

		postId, err := e.Attribute("rel")
		if err != nil {
			loopErr = err
			break
		}

		sourcePath := "/?p=" + *postId

		href, err := e.Attribute("href")
		if err != nil {
			loopErr = err
			break
		}

		slug := helper.GetSlug(*href)

		results = append(results, v1Model.SeriesList{
			Title:      title,
			Slug:       slug,
			SourcePath: sourcePath,
		})
	}

	if loopErr != nil {
		return nil, loopErr
	}

	return results, nil
}
