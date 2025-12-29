package domain

// DashboardStats contains aggregated statistics for the dashboard.
type DashboardStats struct {
	TotalCoins           int64          `json:"total_coins"`
	TotalValue           float64        `json:"total_value"`
	AverageValue         float64        `json:"average_value"`
	TopValuableCoins     []Coin         `json:"top_valuable_coins"`
	RecentCoins          []Coin         `json:"recent_coins"`
	MaterialDistribution map[string]int `json:"material_distribution"`
	GradeDistribution    map[string]int `json:"grade_distribution"`
	ValueDistribution    map[string]int `json:"value_distribution"`
	CountryDistribution  map[string]int `json:"country_distribution"`
	CenturyDistribution  map[string]int `json:"century_distribution"`
	DecadeDistribution   map[string]int `json:"decade_distribution"`
	OldestCoin           *Coin          `json:"oldest_coin"`
	OldestHighGradeCoin  *Coin          `json:"oldest_high_grade_coin"`
	RarestCoins          []Coin         `json:"rarest_coins"`
	GroupDistribution    map[string]int `json:"group_distribution"`
	TotalSilverWeight    float64        `json:"total_silver_weight"`
	TotalGoldWeight      float64        `json:"total_gold_weight"`
	TotalSilverValue     float64        `json:"total_silver_value"`
	TotalGoldValue       float64        `json:"total_gold_value"`
	HeaviestCoin         *Coin          `json:"heaviest_coin"`
	SmallestCoin         *Coin          `json:"smallest_coin"`
	RandomCoin           *Coin          `json:"random_coin"`
	AllCoins             []Coin         `json:"all_coins"`
	GroupStats           []GroupStat    `json:"group_stats"`
}

type GroupStat struct {
	GroupID   int     `json:"group_id"`
	GroupName string  `json:"group_name"`
	Count     int64   `json:"count"`
	MinVal    float64 `json:"min_value"`
	MaxVal    float64 `json:"max_value"`
}
