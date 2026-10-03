package gilliesettings

import repo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/gillie_settings"

type SettingsData struct {
	BaseURL    string `json:"base_url"`
	APIKey     string `json:"api_key"`
	Configured bool   `json:"configured"`
	UpdatedBy  string `json:"updated_by"`
}

func toSettingsData(stored repo.Settings) SettingsData {
	return SettingsData{
		BaseURL:    stored.BaseURL,
		APIKey:     masked(stored.APIKeyLast4),
		Configured: stored.BaseURL != "" && stored.APIKeyEncrypted != "",
		UpdatedBy:  stored.UpdatedBy,
	}
}

func masked(last4 string) string {
	if last4 == "" {
		return ""
	}
	return "••••••••" + last4
}
