//go:generate go run github.com/steebchen/prisma-client-go generate --schema=../../schema.prisma
package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/getsentry/sentry-go"
	"github.com/opensearch-project/opensearch-go/v2"

	_ "fourleaves.studio/manga-scraper/docs"
	"fourleaves.studio/manga-scraper/internal/config"
	"fourleaves.studio/manga-scraper/internal/database/prisma"
	server "fourleaves.studio/manga-scraper/internal/rest"
)

// @title						Manga Scraper API
// @version					1.0
// @description				This is a Manga Scraper API server.
// @termsOfService				https://manga-scraper.hostinger.fourleaves.studio/terms
// @contact.name				API Support
// @contact.url				https://manga-scraper.hostinger.fourleaves.studio/support
// @contact.email				admin@fourleaves.studio
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @securitydefinitions.apikey	TokenAuth
// @in							header
// @name						Authorization
// @tokenUrl					https://manga-reader.fourleaves.studio/sign-in
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

	kafkaClient, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": envConfig.KafkaURL,
	})
	if err != nil {
		log.Fatal("[main] failed to create kafka client: ", err)
	}

	srv := server.NewRESTServer(envConfig, dbClient, esClient, kafkaClient)
	errC, err := srv.StartServer()
	if err != nil {
		log.Fatal("[main] couldn't run: ", err)
	}

	if err := <-errC; err != nil {
		log.Fatal("[main] error while running: ", err)
	}
}
