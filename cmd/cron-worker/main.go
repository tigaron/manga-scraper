package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/opensearch-project/opensearch-go/v2"
	"go.uber.org/zap"

	"fourleaves.studio/manga-scraper/internal/config"
	"fourleaves.studio/manga-scraper/internal/cron"
	"fourleaves.studio/manga-scraper/internal/database/prisma"
	"fourleaves.studio/manga-scraper/internal/elasticsearch"
	kafkaDomain "fourleaves.studio/manga-scraper/internal/kafka"
	"fourleaves.studio/manga-scraper/internal/service"
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

	dbClient := prisma.NewClient(prisma.WithDatasourceURL(envConfig.DBURL))
	if err := dbClient.Connect(); err != nil {
		log.Fatal("[main] failed to connect to database: ", err)
	}
	defer func() {
		if err := dbClient.Disconnect(); err != nil {
			log.Fatal("[main] failed to disconnect from database: ", err)
		}
	}()

	esClient, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // nolint:gosec
		},
		Addresses: []string{envConfig.SearchURL},
	})
	if err != nil {
		log.Fatal("[main] failed to create elasticsearch client: ", err)
	}

	kafkaClient, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": envConfig.KafkaURL,
	})
	if err != nil {
		log.Fatal("[main] failed to create kafka client: ", err)
	}
	defer kafkaClient.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("[main] failed to create logger: ", err)
	}

	cronRepo := prisma.NewCronJobRepo(dbClient)
	providerRepo := prisma.NewProviderRepo(dbClient)
	seriesRepo := prisma.NewSeriesRepo(dbClient)

	seriesSearch := elasticsearch.NewSeriesSearchRepository(esClient)

	scraperRepo := prisma.NewScraperRepo(dbClient)
	scaperMessageBroker := kafkaDomain.NewScraperMessageBroker(kafkaClient)
	scraperService := service.NewScraperCronService(scraperRepo, scaperMessageBroker, logger)

	cronWorker := cron.NewCron(providerRepo, seriesRepo, cronRepo, scraperService, seriesSearch, logger)

	errC, err := cronWorker.StartServer()
	if err != nil {
		log.Fatal("[main] couldn't run: ", err)
	}

	if err := <-errC; err != nil {
		log.Fatal("[main] error while running: ", err)
	}
}
