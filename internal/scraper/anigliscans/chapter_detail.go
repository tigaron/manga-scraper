package anigliscans

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
)

func ScrapeChapterDetail(ctx context.Context, browserUrl, chapterUrl string) (internal.ChapterDetailResult, error) {
	l, err := launcher.NewManaged(browserUrl)
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	l.Leakless(true)
	l.Headless(true)

	lC, err := l.Client()
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	browser := rod.New().Client(lC)
	err = browser.Connect()
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	defer browser.MustClose()

	pg, err := browser.Page(proto.TargetCreateTarget{URL: chapterUrl})
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	page := pg.Context(ctx)

	elFT, err := page.Element("h1.entry-title")
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	tFT, err := elFT.Text()
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	fullTitle := strings.TrimSpace(tFT)

	elHR, err := page.Element("link[rel='shortlink']")
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	href, err := elHR.Attribute("href")
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	postId := helper.GetPostId(*href)
	sourcePath := "/?p=" + postId

	elTS, err := page.ElementR("script", "ts_reader.run")
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	tTS, err := elTS.Text()
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	var tsReader internal.TSReaderScript

	scriptRegex := regexp.MustCompile(`^\s*ts_reader.run\((.*)\);`)
	script := scriptRegex.FindStringSubmatch(tTS)[1]

	err = json.Unmarshal([]byte(script), &tsReader)
	if err != nil {
		return internal.ChapterDetailResult{}, err
	}

	var contentPaths []string

	images := tsReader.Sources[0].Images
	for i := range images {
		if images[i] == "" {
			continue
		}

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
		return internal.ChapterDetailResult{}, err
	}

	nextHref := tsReader.NextURL
	nextSlug := helper.GetSlug(nextHref)

	prevHref := tsReader.PrevURL
	prevSlug := helper.GetSlug(prevHref)

	var nextPath string

	if nextHref != "" {
		err := page.Navigate(nextHref)
		if err != nil {
			return internal.ChapterDetailResult{}, err
		}

		elHRN, err := page.Element("link[rel='shortlink']")
		if err != nil {
			return internal.ChapterDetailResult{}, err
		}

		hrefN, err := elHRN.Attribute("href")
		if err != nil {
			return internal.ChapterDetailResult{}, err
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
			return internal.ChapterDetailResult{}, err
		}

		elHRP, err := page.Element("link[rel='shortlink']")
		if err != nil {
			return internal.ChapterDetailResult{}, err
		}

		hrefP, err := elHRP.Attribute("href")
		if err != nil {
			return internal.ChapterDetailResult{}, err
		}

		postIdP := helper.GetPostId(*hrefP)
		prevPath = "/?p=" + postIdP
	} else {
		prevPath = ""
	}

	result := internal.ChapterDetailResult{
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