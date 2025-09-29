package model

// Common structures used across the project (shared by parser and timeline)

type RateLimitDetail struct {
	UsedPercent float64 `json:"used_percent"`
}

type RateLimits struct {

	Primary              RateLimitDetail `json:"primary"`
	Secondary            RateLimitDetail `json:"secondary"`
	PrimaryUsedPercent   float64         `json:"primary_used_percent"`
	SecondaryUsedPercent float64         `json:"secondary_used_percent"`

}

type TokenCount struct {
	Info struct {
		TotalTokenUsage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"total_token_usage"`
		LastTokenUsage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"last_token_usage"`
	} `json:"info"`
	RateLimits RateLimits `json:"rate_limits"`
}

type LogLine struct {
	Timestamp string     `json:"timestamp"`
	Type      string     `json:"type"`
	Payload   TokenCount `json:"payload"`
}

// Timeline entry with already extracted values
type TimelineEntry struct {
	Timestamp string
	Primary   int
	Secondary int
}
