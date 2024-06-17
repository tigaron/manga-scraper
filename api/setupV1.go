package api

import (
	"fourleaves.studio/manga-scraper/api/middlewares"
	"github.com/labstack/echo/v4"
)

func (s *RESTServer) setupV1Router(v1Api *echo.Group) {
	providers := v1Api.Group("/providers")
	providers.GET("", s.v1.GetProvidersList)
	providers.POST("", s.v1.PostProvider, middlewares.IsAdmin(s.config.AdminSub))
	providers.GET("/:provider_slug", s.v1.GetProvider)
	providers.GET("/:provider_slug/_bc", s.v1.GetProviderBreadcrumbs)
	providers.PUT("/:provider_slug", s.v1.PutProvider, middlewares.IsAdmin(s.config.AdminSub))

	scrapeRequests := v1Api.Group("/scrape-requests", middlewares.IsAdmin(s.config.AdminSub))
	scrapeRequests.POST("/series/list", s.v1.PostScrapeSeriesList)
	scrapeRequests.PUT("/series/detail", s.v1.PutScrapeSeriesDetail)
	scrapeRequests.POST("/chapters/list", s.v1.PostScrapeChapterList)
	scrapeRequests.PUT("/chapters/detail", s.v1.PutScrapeChapterDetail)

	series := v1Api.Group("/series")
	series.GET("/:provider_slug", s.v1.GetSeriesListPaginated)
	series.GET("/:provider_slug/_all", s.v1.GetSeriesListAll)
	series.GET("/:provider_slug/:series_slug", s.v1.GetSeries)
	series.GET("/:provider_slug/:series_slug/_bc", s.v1.GetSeriesBreadcrumbs)
	series.PUT("/:provider_slug/:series_slug/_chc", s.v1.PutSeriesChaptersCount, middlewares.IsAdmin(s.config.AdminSub))
	series.PUT("/:provider_slug/:series_slug/_lch", s.v1.PutSeriesLastChapter, middlewares.IsAdmin(s.config.AdminSub))

	chapters := v1Api.Group("/chapters")
	chapters.GET("/:provider_slug/:series_slug", s.v1.GetChapterListPaginated)
	chapters.GET("/:provider_slug/:series_slug/_list", s.v1.GetChapterList)
	chapters.GET("/:provider_slug/:series_slug/_all", s.v1.GetChapterListAll)
	chapters.GET("/:provider_slug/:series_slug/:chapter_slug", s.v1.GetChapter)
	chapters.GET("/:provider_slug/:series_slug/:chapter_slug/_bc", s.v1.GetChapterBreadcrumbs)

	search := v1Api.Group("/search")
	search.GET("", s.v1.GetSearch)
	search.PUT("", s.v1.PutSearch, middlewares.IsAdmin(s.config.AdminSub))
}
