package mangagalaxy

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	v1Model "fourleaves.studio/manga-scraper/api/models/v1"
	"fourleaves.studio/manga-scraper/internal/scraper/helper"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func ScrapeSeriesDetail(ctx context.Context, browserUrl, seriesUrl string) (v1Model.SeriesDetail, error) {
	l, err := launcher.NewManaged(browserUrl)
	if err != nil {
		return v1Model.SeriesDetail{}, err
	}

	l.Leakless(true)
	l.Headless(true)

	lC, err := l.Client()
	if err != nil {
		return v1Model.SeriesDetail{}, err
	}

	browser := rod.New().Client(lC)
	err = browser.Connect()
	if err != nil {
		return v1Model.SeriesDetail{}, err
	}

	defer browser.MustClose()

	pg, err := browser.Page(proto.TargetCreateTarget{URL: seriesUrl})
	if err != nil {
		return v1Model.SeriesDetail{}, err
	}

	page := pg.Context(ctx)

	elTH, err := page.Element("div.thumb > img")
	if err != nil {
		return v1Model.SeriesDetail{}, err
	}

	thumbnailUrl, err := elTH.Attribute("src")
	if err != nil {
		return v1Model.SeriesDetail{}, err
	}

	if !strings.HasPrefix(*thumbnailUrl, "http") {
		thumbnailUrl, err = elTH.Attribute("data-src")
		if err != nil {
			return v1Model.SeriesDetail{}, err
		}
	}

	// ex of wrong url
	// "https://i0.wp.com/mangagalaxy.me/wp-content/uploads/Guy-in-center-with-crystal-3957.jpg"
	// "https://i1.wp.com/mangagalaxy.me/wp-content/uploads/Guy-in-center-with-crystal-3957.jpg"
	// "https://i3.wp.com/mangagalaxy.me/wp-content/uploads/Guy-in-center-with-crystal-3957.jpg"
	// ex of correct url
	// "https://mangagalaxy.me/wp-content/uploads/Guy-in-center-with-crystal-3957.jpg"
	providerHost := strings.Split(seriesUrl, "/")[2]
	thumbnailHost := strings.Split(*thumbnailUrl, "/")[2]
	if providerHost != thumbnailHost {
		newThumbUrl := strings.Join(strings.Split(*thumbnailUrl, thumbnailHost+"/"), "")
		thumbnailUrl = &newThumbUrl
	}

	var synopsisArr []string

	elS, err := page.Element("div.entry-content")
	if err != nil {
		return v1Model.SeriesDetail{}, err
	}

	elSP, err := elS.Elements("p")
	if err != nil {
		return v1Model.SeriesDetail{}, err
	}

	if len(elSP) == 0 {
		elD, err := elS.Element(`div[class^="contents"]`)
		if err != nil {
			return v1Model.SeriesDetail{}, err
		}

		elDt, err := elD.Elements("div")
		if err != nil {
			return v1Model.SeriesDetail{}, err
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
		return v1Model.SeriesDetail{}, err
	}

	for _, e := range elG {
		genre, _ := e.Text()
		genreArr = append(genreArr, genre)
	}

	genres, err := json.Marshal(helper.RemoveDuplicate(genreArr))
	if err != nil {
		return v1Model.SeriesDetail{}, err
	}

	result := v1Model.SeriesDetail{
		ThumbnailURL: *thumbnailUrl,
		Synopsis:     synopsis,
		Genres:       genres,
	}

	return result, nil
}
