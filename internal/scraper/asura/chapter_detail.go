package asura

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

func ScrapeChapterDetail(ctx context.Context, browserUrl, chapterUrl string) (v1Model.ChapterDetail, error) {
	l, err := launcher.NewManaged(browserUrl)
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	l.Leakless(true)
	l.Headless(true)

	lC, err := l.Client()
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	browser := rod.New().Client(lC)
	err = browser.Connect()
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	defer browser.MustClose()

	pg, err := browser.Page(proto.TargetCreateTarget{URL: chapterUrl})
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	page := pg.Context(ctx)

	elFT, err := page.Element("h1.entry-title")
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	tFT, err := elFT.Text()
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	fullTitle := strings.TrimSpace(tFT)

	elHR, err := page.Element("link[rel='shortlink']")
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	href, err := elHR.Attribute("href")
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	postId := helper.GetPostId(*href)
	sourcePath := "/?p=" + postId

	elTS, err := page.ElementR("script", "ts_reader.run")
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	tTS, err := elTS.Text()
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	var tsReader v1Model.TSReaderScript

	scriptRegex := regexp.MustCompile(`^\s*ts_reader.run\((.*)\);`)
	script := scriptRegex.FindStringSubmatch(tTS)[1]

	err = json.Unmarshal([]byte(script), &tsReader)
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	var contentPaths []string

	images := tsReader.Sources[0].Images
	for i := range images {
		imgSplit := strings.Split(images[i], "/")
		imgPath := strings.Join(imgSplit[3:], "/")

		// only append if imgSplit has more than 3 elements, ex:
		// https://i.imgur.com/5yKo93E.jpg
		// https://asuratoon.com/wp-content/uploads/custom-upload/96904/92/00 kopya.jpg
		if len(imgSplit) > 3 {
			contentPaths = append(contentPaths, "/"+imgPath)
		}
	}

	contentPathsJson, err := json.Marshal(helper.RemoveDuplicate(contentPaths))
	if err != nil {
		return v1Model.ChapterDetail{}, err
	}

	nextHref := tsReader.NextURL
	nextSlug := helper.GetSlug(nextHref)

	prevHref := tsReader.PrevURL
	prevSlug := helper.GetSlug(prevHref)

	var nextPath string

	if nextHref != "" {
		err := page.Navigate(nextHref)
		if err != nil {
			return v1Model.ChapterDetail{}, err
		}

		elHRN, err := page.Element("link[rel='shortlink']")
		if err != nil {
			return v1Model.ChapterDetail{}, err
		}

		hrefN, err := elHRN.Attribute("href")
		if err != nil {
			return v1Model.ChapterDetail{}, err
		}

		postIdN := helper.GetPostId(*hrefN)
		nextPath = "/?p=" + postIdN
	} else {
		nextPath = ""
	}

	var prevPath string

	if prevHref != "" {
		err := page.Navigate(prevHref)
		if err != nil {
			return v1Model.ChapterDetail{}, err
		}

		elHRP, err := page.Element("link[rel='shortlink']")
		if err != nil {
			return v1Model.ChapterDetail{}, err
		}

		hrefP, err := elHRP.Attribute("href")
		if err != nil {
			return v1Model.ChapterDetail{}, err
		}

		postIdP := helper.GetPostId(*hrefP)
		prevPath = "/?p=" + postIdP
	} else {
		prevPath = ""
	}

	result := v1Model.ChapterDetail{
		FullTitle:    fullTitle,
		SourcePath:   sourcePath,
		ContentPaths: contentPathsJson,
		NextPath:     nextPath,
		NextSlug:     nextSlug,
		PrevPath:     prevPath,
		PrevSlug:     prevSlug,
	}

	return result, nil
}
