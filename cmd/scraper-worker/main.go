package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/getsentry/sentry-go"
	"github.com/opensearch-project/opensearch-go/v2"
	"go.uber.org/zap"

	"fourleaves.studio/manga-scraper/internal/config"
	"fourleaves.studio/manga-scraper/internal/database/prisma"
	"fourleaves.studio/manga-scraper/internal/elasticsearch"
	"fourleaves.studio/manga-scraper/internal/scraper"
)

func main() {
	// Set local timezone to Asia/Singapore
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		log.Fatal("[main] failed to load location: ", err)
	}

	time.Local = loc

	// Load config from .env file
	envConfig, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal("[main] failed to load config: ", err)
	}

	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           envConfig.SentryDSN,
		Environment:   envConfig.ENV,
		Release:       envConfig.Version,
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate:   1.0,
		IgnoreTransactions: []string{"/health", "/swagger"},
	}); err != nil {
		log.Fatal("[main] failed to initialize sentry: ", err)
	}

	dbClient := prisma.NewClient(prisma.WithDatasourceURL(envConfig.DBURL))
	if err := dbClient.Connect(); err != nil {
		log.Fatal("[main] failed to connect to database: ", err)
	}

	esClient, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // nolint:gosec
		},
		Addresses: []string{envConfig.SearchURL},
	})
	if err != nil {
		log.Fatal("[main] failed to create elasticsearch client: ", err)
	}

	kafkaClient, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  envConfig.KafkaURL,
		"sasl.mechanism":     "SCRAM-SHA-256",
		"security.protocol":  "SASL_SSL",
		"sasl.username":      envConfig.KafkaUsername,
		"sasl.password":      envConfig.KafkaPassword,
		"group.id":           "scraper-worker",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		log.Fatal("[main] failed to create kafka client: ", err)
	}

	if err := kafkaClient.Subscribe("scrape-request", nil); err != nil {
		log.Fatal("[main] failed to subscribe to kafka topic: ", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("[main] failed to create logger: ", err)
	}

	seriesRepo := prisma.NewSeriesRepo(dbClient)
	chapterRepo := prisma.NewChapterRepo(dbClient)
	scraperRepo := prisma.NewScraperRepo(dbClient)

	seriesSearch := elasticsearch.NewSeriesSearchRepository(esClient)

	scraperService := scraper.NewScraper(scraperRepo, seriesRepo, seriesSearch, chapterRepo, kafkaClient, logger, envConfig.RodURL)

	errC, err := scraperService.StartServer()
	if err != nil {
		log.Fatal("[main] couldn't run: ", err)
	}

	if err := <-errC; err != nil {
		log.Fatal("[main] error while running: ", err)
	}
}
