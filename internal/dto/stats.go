package dto

// StatsQuery contains filters for stats summary and export.
type StatsQuery struct {
	Period    string  `form:"period" binding:"required" validate:"required,oneof=day week month year"`
	Value     string  `form:"value" binding:"required" validate:"required,min=1,max=30"`
	AccountID *string `form:"account_id,omitempty"`
}

// TimeseriesPoint contains one bucket in the stats chart.
type TimeseriesPoint struct {
	Label   string  `json:"label"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

// CategoryBreakdownItem contains per-category aggregation.
type CategoryBreakdownItem struct {
	Category    string  `json:"category"`
	Subcategory string  `json:"subcategory"`
	Amount      float64 `json:"amount"`
	Percentage  float64 `json:"percentage"`
}

// StatsSummaryResponse contains overall stats response payload.
type StatsSummaryResponse struct {
	TotalIncome              float64                 `json:"total_income"`
	TotalExpense             float64                 `json:"total_expense"`
	Net                      float64                 `json:"net"`
	CategoryBreakdownExpense []CategoryBreakdownItem `json:"category_breakdown_expense"`
	CategoryBreakdownIncome  []CategoryBreakdownItem `json:"category_breakdown_income"`
	Timeseries               []TimeseriesPoint       `json:"timeseries"`
}
