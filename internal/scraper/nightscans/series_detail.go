package nightscans

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"fourleaves.studio/manga-scraper/internal"
	"fourleaves.studio/manga-scraper/internal/scraper/helper"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
)

func ScrapeSeriesDetail(ctx context.Context, browserURL, seriesURL string, logger *zap.Logger) (internal.SeriesDetailResult, error) {
	l, err := launcher.NewManaged(browserURL)
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	l.Leakless(true)
	l.Headless(true)

	lC, err := l.Client()
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	browser := rod.New().Client(lC)
	err = browser.Connect()
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	defer browser.MustClose()

	pg, err := browser.Page(proto.TargetCreateTarget{URL: seriesURL})
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	page := pg.Context(ctx)

	elTH, err := page.Element("div.thumb > img")
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	thumbnailURL, err := elTH.Attribute("src")
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	if !strings.HasPrefix(*thumbnailURL, "http") {
		thumbnailURL, err = elTH.Attribute("data-lazy-src")
		if err != nil {
			return internal.SeriesDetailResult{}, err
		}
	}

	var synopsisArr []string

	elS, err := page.Element("div.entry-content")
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	elSP, err := elS.Elements("p")
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	if len(elSP) == 0 {
		elD, err := elS.Element(`div[class^="contents"]`)
		if err != nil {
			return internal.SeriesDetailResult{}, err
		}

		elDt, err := elD.Elements("div")
		if err != nil {
			return internal.SeriesDetailResult{}, err
		}

		for _, e := range elDt {
			text, _ := e.Text()
			if text == "&nbsp;" || text == "\u00a0" || text == "" {
				continue
			}

			synopsisArr = append(synopsisArr, text)
		}
	} else {
		for _, e := range elSP {
			text, _ := e.Text()
			if text == "&nbsp;" || text == "\u00a0" || text == "" {
				continue
			}

			synopsisArr = append(synopsisArr, text)
		}
	}

	synopsisRegex := regexp.MustCompile(`\n`)
	synopsis := synopsisRegex.ReplaceAllString(strings.Join(helper.RemoveDuplicate(synopsisArr), "<br />"), "<br />")

	var genreArr []string

	elG, err := page.Elements("span.mgen > a")
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	for _, e := range elG {
		genre, _ := e.Text()
		genreArr = append(genreArr, genre)
	}

	genres, err := json.Marshal(helper.RemoveDuplicate(genreArr))
	if err != nil {
		return internal.SeriesDetailResult{}, err
	}

	result := internal.SeriesDetailResult{
		ThumbnailURL: *thumbnailURL,
		Synopsis:     synopsis,
		Genres:       genres,
	}

	logger.Debug("Scraped series detail", zap.Any("result", result))

	return result, nil
}
