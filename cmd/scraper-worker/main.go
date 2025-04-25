package main

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"

	"fourleaves.studio/manga-scraper/internal/config"
	"fourleaves.studio/manga-scraper/internal/database/prisma"
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

	dbClient := prisma.NewClient(prisma.WithDatasourceURL(envConfig.DBURL))
	if err := dbClient.Connect(); err != nil {
		log.Fatal("[main] failed to connect to database: ", err)
	}

	kafkaClient, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  envConfig.KafkaURL,
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

	scraperService := scraper.NewScraper(scraperRepo, seriesRepo, chapterRepo, kafkaClient, logger, envConfig.RodURL)

	errC, err := scraperService.StartServer()
	if err != nil {
		log.Fatal("[main] couldn't run: ", err)
	}

	if err := <-errC; err != nil {
		log.Fatal("[main] error while running: ", err)
	}
}
