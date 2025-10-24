package dto

type LibListItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GitURL      string `json:"git_url"`
}

type CreateLibRequest struct {
	GitURL string `json:"git_url"`
}
