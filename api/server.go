package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	v1Handler "fourleaves.studio/manga-scraper/api/handlers/v1"
	"fourleaves.studio/manga-scraper/api/middlewares"
	"fourleaves.studio/manga-scraper/internal/config"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"fourleaves.studio/manga-scraper/internal/database/redis"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type RESTServer struct {
	router *echo.Echo
	config *config.Config
	v1     *v1Handler.Handler
}

func NewRESTServer(config *config.Config, db *db.PrismaClient, redis *redis.RedisClient) *RESTServer {
	v1 := v1Handler.NewV1Handler(config, db, redis)

	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Logger.SetLevel(log.INFO)

	app.Use(middlewares.SentryMiddleware())

	app.Use(middlewares.TimeoutMiddleware(30 * time.Second))

	clerk.SetKey(config.ClerkSecretKey)

	app.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodOptions},
		}),
	)

	app.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	server := &RESTServer{
		router: app,
		config: config,
		v1:     v1,
	}

	app.GET("/health", server.v1.GetHealthCheck)

	app.GET("/swagger/*", echoSwagger.WrapHandler)

	v1Api := app.Group("/api/v1")

	server.setupV1Router(v1Api)

	return server
}

func (s *RESTServer) StartServer(port string) {
	// Setup context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := s.router.Start(port); err != nil && err != http.ErrServerClosed {
			s.router.Logger.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.router.Shutdown(ctx); err != nil {
		s.router.Logger.Fatal(err)
	}
}

func (s *RESTServer) setupV1Router(v1Api *echo.Group) {
	providers := v1Api.Group("/providers")
	providers.GET("", s.v1.GetProvidersList)
	providers.POST("", s.v1.PostProvider)
	providers.GET("/:provider_slug", s.v1.GetProvider)
	providers.PUT("/:provider_slug", s.v1.PutProvider)

	scrapeRequests := v1Api.Group("/scrape-requests", middlewares.WithHeaderAuth)
	scrapeRequests.POST("/series/list", s.v1.PostScrapeSeriesList)
	scrapeRequests.PUT("/series/detail", s.v1.PutScrapeSeriesDetail)
	scrapeRequests.POST("/chapters/list", s.v1.PostScrapeChapterList)
	scrapeRequests.PUT("/chapters/detail", s.v1.PutScrapeChapterDetail)

	series := v1Api.Group("/series")
	series.GET("/:provider_slug", s.v1.GetSeriesListPaginated)
	series.GET("/:provider_slug/_all", s.v1.GetSeriesListAll)
	series.GET("/:provider_slug/:series_slug", s.v1.GetSeries)

	chapters := v1Api.Group("/chapters", middlewares.WithHeaderAuth)
	chapters.GET("/:provider_slug/:series_slug", s.v1.GetChapterListPaginated)
	chapters.GET("/:provider_slug/:series_slug/_list", s.v1.GetChapterList)
	chapters.GET("/:provider_slug/:series_slug/_all", s.v1.GetChapterListAll)
	chapters.GET("/:provider_slug/:series_slug/:chapter_slug", s.v1.GetChapter)

	search := v1Api.Group("/search")
	search.GET("", s.v1.GetSearch)
	search.PUT("", s.v1.PutSearch)
}
