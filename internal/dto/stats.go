package dto

// StatsQuery contains filters for stats summary and export.
type StatsQuery struct {
	Period    string  `form:"period" binding:"required" validate:"required,oneof=day week month year"`
	Value     string  `form:"value" binding:"required" validate:"required,min=1,max=30"`
	AccountID *string `form:"account_id,omitempty"`
}

// CompareQuery contains filters for stats comparison.
type CompareQuery struct {
	Period    string `form:"period" binding:"required" validate:"required,oneof=day week month year"`
	Value1    string `form:"value1" binding:"required" validate:"required,min=1,max=30"`
	Value2    string `form:"value2" binding:"required" validate:"required,min=1,max=30"`
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

// ComparePeriodData contains stats for one comparison period.
type ComparePeriodData struct {
	Label        string                  `json:"label"`
	TotalIncome  float64                 `json:"total_income"`
	TotalExpense float64                 `json:"total_expense"`
	Net          float64                 `json:"net"`
	TopExpense   []CategoryBreakdownItem `json:"top_expense"`
	TopIncome    []CategoryBreakdownItem `json:"top_income"`
}

// CompareResponse contains the comparison of two periods.
type CompareResponse struct {
	Period1          ComparePeriodData `json:"period1"`
	Period2          ComparePeriodData `json:"period2"`
	IncomeChange     float64           `json:"income_change_pct"`
	ExpenseChange    float64           `json:"expense_change_pct"`
	NetChange        float64           `json:"net_change"`
}
