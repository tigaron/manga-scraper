package asura

import (
	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	"fourleaves.studio/manga-scraper/internal/scraper/helper"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func ScrapeChapterList(seriesUrl string) ([]v1Model.ChapterList, error) {
	debugURL, err := launcher.New().Leakless(true).Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(debugURL)
	err = browser.Connect()
	if err != nil {
		return nil, err
	}

	defer browser.MustClose()

	page, err := browser.Page(proto.TargetCreateTarget{URL: seriesUrl})
	if err != nil {
		return nil, err
	}

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
