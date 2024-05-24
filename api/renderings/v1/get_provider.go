package v1Response

import db "fourleaves.studio/manga-scraper/internal/database/prisma"

type GetProviderData struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	BaseURL string `json:"baseURL"`
}

func NewGetProviderData(provider *db.ProviderModel) GetProviderData {
	return GetProviderData{
		Name:    provider.Name,
		Slug:    provider.Slug,
		BaseURL: provider.Scheme + provider.Host,
	}
}

func NewGetProvidersListData(providers []db.ProviderModel) []GetProviderData {
	result := make([]GetProviderData, 0, len(providers))
	for _, provider := range providers {
		result = append(result, NewGetProviderData(&provider))
	}

	return result
}
