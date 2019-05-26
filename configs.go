package main

const (
	// APIEndpoint bot%s - token
	// APIEndpoint %s - method
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
)

// UpdateConfig Config for updates
type UpdateConfig struct {
	offset  int
	limit   int
	timeout int
}
