package asura

import (
	"context"

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

	var loopErr error

	for _, e := range elAs {
		href, err := e.Attribute("href")
		if err != nil {
			loopErr = err
			break
		}

		slug := helper.GetSlug(*href)

		elT, err := e.Element("span.chapternum")
		if err != nil {
			loopErr = err
			break
		}

		tT, err := elT.Text()
		if err != nil {
			loopErr = err
			break
		}

		title := helper.GetChapterTitle(tT)

		li, err := e.Parents("li")
		if err != nil {
			loopErr = err
			break
		}

		dataNum, err := li.First().Attribute("data-num")
		if err != nil {
			loopErr = err
			break
		}

		chapterNumber := helper.GetChapterNumber(*dataNum)

		results = append(results, v1Model.ChapterList{
			ShortTitle: title,
			Slug:       slug,
			Number:     chapterNumber,
			Href:       *href,
		})
	}

	if loopErr != nil {
		return nil, loopErr
	}

	return results, nil
}
