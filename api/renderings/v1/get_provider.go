package v1Response

import db "fourleaves.studio/manga-scraper/internal/database/prisma"

type ProviderData struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	BaseURL string `json:"baseURL"`
}

func NewProviderData(provider *db.ProviderModel) ProviderData {
	return ProviderData{
		Name:    provider.Name,
		Slug:    provider.Slug,
		BaseURL: provider.Scheme + provider.Host,
	}
}

func NewProvidersListData(providers []db.ProviderModel) []ProviderData {
	result := make([]ProviderData, 0, len(providers))
	for _, provider := range providers {
		result = append(result, NewProviderData(&provider))
	}

	return result
}
