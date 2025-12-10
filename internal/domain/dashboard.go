package domain

// DashboardStats contains aggregated statistics for the dashboard.
type DashboardStats struct {
	TotalCoins           int64          `json:"total_coins"`
	TotalValue           float64        `json:"total_value"`
	TopValuableCoins     []Coin         `json:"top_valuable_coins"`
	RecentCoins          []Coin         `json:"recent_coins"`
	ValueDistribution    map[string]int `json:"value_distribution"`
	MaterialDistribution map[string]int `json:"material_distribution"`
	GradeDistribution    map[string]int `json:"grade_distribution"`
}
