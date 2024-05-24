//go:generate go run github.com/steebchen/prisma-client-go generate --schema=../../schema.prisma
package main

import (
	"log"
	"time"

	"fourleaves.studio/manga-scraper/api"
	_ "fourleaves.studio/manga-scraper/docs"
	"fourleaves.studio/manga-scraper/internal/config"
	db "fourleaves.studio/manga-scraper/internal/database/prisma"
	"fourleaves.studio/manga-scraper/internal/database/redis"
	"github.com/getsentry/sentry-go"
)

var (
	envConfig   *config.Config
	dbClient    *db.PrismaClient
	redisClient *redis.RedisClient
)

func init() {
	// Set local timezone to Asia/Singapore
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		log.Fatal("[init] failed to load location: ", err)
	}

	time.Local = loc

	// Load config from .env file
	envConfig, err = config.LoadConfig(".env")
	if err != nil {
		log.Fatal("[init] failed to load config: ", err)
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
		log.Fatal("[init] failed to initialize sentry: ", err)
	}

	dbClient = db.NewClient(db.WithDatasourceURL(envConfig.DBURL))
	if err := dbClient.Connect(); err != nil {
		log.Fatal("[init] failed to connect to database: ", err)
	}

	redisClient, err = redis.NewClient(envConfig.RedisURL)
	if err != nil {
		log.Fatal("[init] failed to connect to redis: ", err)
	}
}

// @title			Manga Scraper API
// @version		1.0
// @description	This is a Manga Scraper API server.
// @termsOfService	https://manga-scraper.fourleaves.studio/terms

// @contact.name	API Support
// @contact.url	https://manga-scraper.fourleaves.studio/support
// @contact.email	admin@fourleaves.studio

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	defer sentry.Flush(2 * time.Second)
	defer func() {
		if err := dbClient.Disconnect(); err != nil {
			log.Fatal("[main] failed to disconnect from database: ", err)
		}
	}()

	server := api.NewRESTServer(envConfig, dbClient, redisClient)

	server.StartServer(":1323")
}
