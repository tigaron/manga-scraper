package internal

type Provider struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
	BaseURL  string `json:"baseURL"`
	ListURL  string `json:"listURL"`
}

type ProviderBC struct {
	Provider Breadcrumb `json:"provider"`
}

type ProviderParams struct {
	Slug     string
	Name     string
	Scheme   string
	Host     string
	ListPath string
	IsActive *bool
}

func (p *ProviderParams) Validate() error {
	if p.Slug == "" {
		return NewErrorf(ErrInvalidInput, "slug is required")
	}

	if p.Name == "" {
		return NewErrorf(ErrInvalidInput, "name is required")
	}

	if p.Scheme == "" {
		return NewErrorf(ErrInvalidInput, "scheme is required")
	}

	if p.Host == "" {
		return NewErrorf(ErrInvalidInput, "host is required")
	}

	if p.ListPath == "" {
		return NewErrorf(ErrInvalidInput, "list path is required")
	}

	return nil
}
