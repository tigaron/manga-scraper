package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fourleaves.studio/manga-scraper/internal/config"
	"fourleaves.studio/manga-scraper/internal/database/prisma"
	"fourleaves.studio/manga-scraper/internal/database/redis"
	"fourleaves.studio/manga-scraper/internal/elasticsearch"
	kafkaDomain "fourleaves.studio/manga-scraper/internal/kafka"
	"fourleaves.studio/manga-scraper/internal/rest/middlewares"
	v1Handler "fourleaves.studio/manga-scraper/internal/rest/v1"
	chapterHandler "fourleaves.studio/manga-scraper/internal/rest/v1/chapters"
	providersHandler "fourleaves.studio/manga-scraper/internal/rest/v1/providers"
	scraperHandler "fourleaves.studio/manga-scraper/internal/rest/v1/scrapers"
	seriesHandler "fourleaves.studio/manga-scraper/internal/rest/v1/series"
	"fourleaves.studio/manga-scraper/internal/service"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/opensearch-project/opensearch-go/v2"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type RESTServer struct {
	router      *echo.Echo
	config      *config.Config
	dbClient    *prisma.PrismaClient
	esClient    *opensearch.Client
	kafkaClient *kafka.Producer
}

func NewRESTServer(config *config.Config, dbClient *prisma.PrismaClient, esClient *opensearch.Client, kafkaClient *kafka.Producer) *RESTServer {
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	mid := middlewares.NewMiddleware(config)

	switch config.ENV {
	case "development":
		router.Logger.SetLevel(log.DEBUG)
		router.Use(mid.TimeoutMiddleware(3 * time.Minute))
	case "production":
		router.Logger.SetLevel(log.INFO)
		router.Use(mid.TimeoutMiddleware(30 * time.Second))
	}

	router.Use(mid.SentryMiddleware())

	clerk.SetKey(config.ClerkSecretKey)

	router.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodOptions},
		}),
	)

	router.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	providerRepo := prisma.NewProviderRepo(dbClient)
	providerCache := redis.NewProviderCache(config.RedisURL, providerRepo, 30*time.Minute, router.Logger)
	providerService := service.NewProviderService(providerCache)
	providersHandler.NewProviderHandler(providerService).Register(router.Group("/api/v1/providers"), mid)

	seriesRepo := prisma.NewSeriesRepo(dbClient)
	seriesCache := redis.NewSeriesCache(config.RedisURL, seriesRepo, 30*time.Minute, router.Logger)
	seriesSearch := elasticsearch.NewSeriesSearchRepository(esClient)
	seriesService := service.NewSeriesService(seriesCache, seriesSearch, router.Logger)
	seriesHandler.NewSeriesHandler(seriesService).Register(router.Group("/api/v1/series"), mid)

	chapterRepo := prisma.NewChapterRepo(dbClient)
	chapterCache := redis.NewChapterCache(config.RedisURL, chapterRepo, 30*time.Minute, router.Logger)
	chapterService := service.NewChapterService(chapterCache)
	chapterHandler.NewChapterHandler(chapterService).Register(router.Group("/api/v1/chapters"))

	scraperRepo := prisma.NewScraperRepo(dbClient)
	scaperMessageBroker := kafkaDomain.NewScraperMessageBroker(kafkaClient)
	scraperService := service.NewScraperService(scraperRepo, scaperMessageBroker, router.Logger)
	scraperHandler.NewScraperHandler(scraperService).Register(router.Group("/api/v1/scrapers"), mid)

	router.GET("/health", v1Handler.GetHealthCheck)

	router.GET("/swagger/*", echoSwagger.WrapHandler)

	return &RESTServer{
		router:      router,
		config:      config,
		dbClient:    dbClient,
		esClient:    esClient,
		kafkaClient: kafkaClient,
	}
}

func (s *RESTServer) StartServer() (<-chan error, error) {
	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		if err := s.router.Start(s.config.Port); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		s.dbClient.Disconnect()
		s.kafkaClient.Close()
		sentry.Flush(2 * time.Second)
		stop()
		cancel()
		close(errC)
	}()

	if err := s.router.Shutdown(ctx); err != nil {
		errC <- err
	}

	return errC, nil
}
