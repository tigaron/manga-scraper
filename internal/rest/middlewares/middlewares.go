package middlewares

import "fourleaves.studio/manga-scraper/internal/config"

type Middleware struct {
	config *config.Config
}

func NewMiddleware(config *config.Config) *Middleware {
	return &Middleware{
		config: config,
	}
}
