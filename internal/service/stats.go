package service

import (
	"context"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/repository"
	"ledgerA/pkg/pdf"
	"sort"
	"time"

	"github.com/google/uuid"
)

type statsService struct {
	txRepo  repository.TransactionRepository
	catRepo repository.CategoryRepository
	subRepo repository.SubcategoryRepository
}

type statsBreakdownKey struct {
	category    string
	subcategory string
}

// NewStatsService creates a new StatsService.
func NewStatsService(txRepo repository.TransactionRepository, catRepo repository.CategoryRepository, subRepo repository.SubcategoryRepository) StatsService {
	return &statsService{txRepo: txRepo, catRepo: catRepo, subRepo: subRepo}
}

func (s *statsService) Summary(ctx context.Context, userID uuid.UUID, query dto.StatsQuery) (*dto.StatsSummaryResponse, error) {
	from, to, err := parseStatsRange(query.Period, query.Value)
	if err != nil {
		return nil, fmt.Errorf("statsService.Summary.Range: %w", err)
	}

	filter := repository.TransactionListFilter{Type: "all", SortBy: "transaction_date", SortDir: "asc", Page: 1, PerPage: 5000}
	if query.AccountID != nil && *query.AccountID != "" && *query.AccountID != "all" {
		parsedID, err := uuid.Parse(*query.AccountID)
		if err != nil {
			return nil, fmt.Errorf("statsService.Summary.ParseAccountID: %w", err)
		}
		filter.AccountID = &parsedID
	}

	items, _, err := s.txRepo.ListByUserID(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("statsService.Summary.List: %w", err)
	}

	categoryNames := map[uuid.UUID]string{}
	subcategoryNames := map[uuid.UUID]string{}
	if s.catRepo != nil {
		categories, _, catErr := s.catRepo.ListByUserID(ctx, userID)
		if catErr == nil {
			for _, category := range categories {
				categoryNames[category.ID] = category.Name
				if s.subRepo == nil {
					continue
				}

				subs, _, subErr := s.subRepo.ListByCategoryID(ctx, userID, category.ID)
				if subErr != nil {
					continue
				}
				for _, sub := range subs {
					subcategoryNames[sub.ID] = sub.Name
				}
			}
		}
	}

	resp := &dto.StatsSummaryResponse{
		CategoryBreakdownExpense: []dto.CategoryBreakdownItem{},
		CategoryBreakdownIncome:  []dto.CategoryBreakdownItem{},
		Timeseries:               []dto.TimeseriesPoint{},
	}

	timeline := map[string]dto.TimeseriesPoint{}
	expenseBreakdown := map[statsBreakdownKey]float64{}
	incomeBreakdown := map[statsBreakdownKey]float64{}

	for _, item := range items {
		txDate := item.TransactionDate
		if txDate.Before(from) || !txDate.Before(to) {
			continue
		}

		label := timelineLabel(query.Period, txDate)
		point := timeline[label]
		point.Label = label

		categoryName := categoryNames[item.CategoryID]
		if categoryName == "" {
			categoryName = "Uncategorized"
		}

		subcategoryName := subcategoryNames[item.SubcategoryID]
		if subcategoryName == "" {
			subcategoryName = categoryName
		}

		key := statsBreakdownKey{category: categoryName, subcategory: subcategoryName}

		if item.Amount >= 0 {
			resp.TotalIncome += item.Amount
			point.Income += item.Amount
			incomeBreakdown[key] += item.Amount
		} else {
			expenseAmount := -item.Amount
			resp.TotalExpense += expenseAmount
			point.Expense += expenseAmount
			expenseBreakdown[key] += expenseAmount
		}

		timeline[label] = point
	}
	resp.Net = resp.TotalIncome - resp.TotalExpense

	labels := make([]string, 0, len(timeline))
	for label := range timeline {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	for _, label := range labels {
		resp.Timeseries = append(resp.Timeseries, timeline[label])
	}

	resp.CategoryBreakdownExpense = buildBreakdownItems(expenseBreakdown, resp.TotalExpense)
	resp.CategoryBreakdownIncome = buildBreakdownItems(incomeBreakdown, resp.TotalIncome)

	return resp, nil
}

func parseStatsRange(period string, value string) (time.Time, time.Time, error) {
	switch period {
	case "day":
		start, err := time.Parse("2006-01-02", value)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return start, start.AddDate(0, 0, 1), nil
	case "week":
		if len(value) != len("2006-W02") {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid week format")
		}
		year, week, err := parseISOWeek(value)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		start := isoWeekStart(year, week)
		return start, start.AddDate(0, 0, 7), nil
	case "month":
		start, err := time.Parse("2006-01", value)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return start, start.AddDate(0, 1, 0), nil
	case "year":
		start, err := time.Parse("2006", value)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return start, start.AddDate(1, 0, 0), nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid period: %s", period)
	}
}

func parseISOWeek(value string) (int, int, error) {
	var year int
	var week int
	if _, err := fmt.Sscanf(value, "%d-W%d", &year, &week); err != nil {
		return 0, 0, err
	}
	if week < 1 || week > 53 {
		return 0, 0, fmt.Errorf("invalid week number: %d", week)
	}
	return year, week, nil
}

func isoWeekStart(year int, week int) time.Time {
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)
	for jan4.Weekday() != time.Monday {
		jan4 = jan4.AddDate(0, 0, -1)
	}
	return jan4.AddDate(0, 0, (week-1)*7)
}

func timelineLabel(period string, date time.Time) string {
	switch period {
	case "year":
		return date.Format("2006-01")
	default:
		return date.Format("2006-01-02")
	}
}

func buildBreakdownItems(values map[statsBreakdownKey]float64, total float64) []dto.CategoryBreakdownItem {
	items := make([]dto.CategoryBreakdownItem, 0, len(values))
	for key, amount := range values {
		percentage := 0.0
		if total > 0 {
			percentage = (amount / total) * 100
		}
		items = append(items, dto.CategoryBreakdownItem{
			Category:    key.category,
			Subcategory: key.subcategory,
			Amount:      amount,
			Percentage:  percentage,
		})
	}

	sort.Slice(items, func(i int, j int) bool {
		return items[i].Amount > items[j].Amount
	})

	return items
}

func (s *statsService) ExportPDF(ctx context.Context, userID uuid.UUID, query dto.StatsQuery) ([]byte, error) {
	summary, err := s.Summary(ctx, userID, query)
	if err != nil {
		return nil, fmt.Errorf("statsService.ExportPDF.Summary: %w", err)
	}

	accountName := "All Accounts"
	if query.AccountID != nil && *query.AccountID != "" && *query.AccountID != "all" {
		// Account name is not critical for PDF; use ID as fallback
		accountName = *query.AccountID
	}

	from, to, err := parseStatsRange(query.Period, query.Value)
	if err != nil {
		return nil, fmt.Errorf("statsService.ExportPDF.Range: %w", err)
	}

	filter := repository.TransactionListFilter{
		Type: "all", SortBy: "transaction_date", SortDir: "asc",
		Page: 1, PerPage: 5000, DateFrom: &from, DateTo: &to,
	}
	if query.AccountID != nil && *query.AccountID != "" && *query.AccountID != "all" {
		parsedID, parseErr := uuid.Parse(*query.AccountID)
		if parseErr == nil {
			filter.AccountID = &parsedID
		}
	}

	items, _, err := s.txRepo.ListByUserID(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("statsService.ExportPDF.List: %w", err)
	}

	categoryNames := s.buildCategoryNames(ctx, userID)
	subcategoryNames := s.buildSubcategoryNames(ctx, userID)

	breakdownRows := make([]pdf.CategoryBreakdownRow, 0)
	for _, item := range summary.CategoryBreakdownExpense {
		breakdownRows = append(breakdownRows, pdf.CategoryBreakdownRow{
			Category: item.Category, Subcategory: item.Subcategory,
			Amount: item.Amount, Percentage: item.Percentage, Type: "expense",
		})
	}
	for _, item := range summary.CategoryBreakdownIncome {
		breakdownRows = append(breakdownRows, pdf.CategoryBreakdownRow{
			Category: item.Category, Subcategory: item.Subcategory,
			Amount: item.Amount, Percentage: item.Percentage, Type: "income",
		})
	}

	txRows := make([]pdf.TransactionRow, 0, len(items))
	for _, item := range items {
		catName := categoryNames[item.CategoryID]
		if catName == "" {
			catName = "Uncategorized"
		}
		subName := subcategoryNames[item.SubcategoryID]
		notes := ""
		if item.Notes != nil {
			notes = *item.Notes
		}
		txRows = append(txRows, pdf.TransactionRow{
			Date:        item.TransactionDate.Format("2006-01-02"),
			Name:        item.Name,
			Category:    catName,
			Subcategory: subName,
			Amount:      item.Amount,
			Notes:       notes,
		})
	}

	periodLabel := fmt.Sprintf("%s: %s", query.Period, query.Value)

	pdfData := pdf.StatsPDFData{
		PeriodLabel:   periodLabel,
		AccountName:   accountName,
		CurrencyCode:  "INR",
		TotalIncome:   summary.TotalIncome,
		TotalExpense:  summary.TotalExpense,
		Net:           summary.Net,
		BreakdownRows: breakdownRows,
		Transactions:  txRows,
	}

	result, err := pdf.GenerateStatsPDF(pdfData)
	if err != nil {
		return nil, fmt.Errorf("statsService.ExportPDF.Generate: %w", err)
	}
	return result, nil
}

func (s *statsService) Compare(ctx context.Context, userID uuid.UUID, query dto.CompareQuery) (*dto.CompareResponse, error) {
	summary1, err := s.Summary(ctx, userID, dto.StatsQuery{
		Period: query.Period, Value: query.Value1, AccountID: query.AccountID,
	})
	if err != nil {
		return nil, fmt.Errorf("statsService.Compare.Period1: %w", err)
	}

	summary2, err := s.Summary(ctx, userID, dto.StatsQuery{
		Period: query.Period, Value: query.Value2, AccountID: query.AccountID,
	})
	if err != nil {
		return nil, fmt.Errorf("statsService.Compare.Period2: %w", err)
	}

	topN := 5
	p1TopExpense := summary1.CategoryBreakdownExpense
	if len(p1TopExpense) > topN {
		p1TopExpense = p1TopExpense[:topN]
	}
	p1TopIncome := summary1.CategoryBreakdownIncome
	if len(p1TopIncome) > topN {
		p1TopIncome = p1TopIncome[:topN]
	}
	p2TopExpense := summary2.CategoryBreakdownExpense
	if len(p2TopExpense) > topN {
		p2TopExpense = p2TopExpense[:topN]
	}
	p2TopIncome := summary2.CategoryBreakdownIncome
	if len(p2TopIncome) > topN {
		p2TopIncome = p2TopIncome[:topN]
	}

	incomeChange := 0.0
	if summary1.TotalIncome > 0 {
		incomeChange = ((summary2.TotalIncome - summary1.TotalIncome) / summary1.TotalIncome) * 100
	}
	expenseChange := 0.0
	if summary1.TotalExpense > 0 {
		expenseChange = ((summary2.TotalExpense - summary1.TotalExpense) / summary1.TotalExpense) * 100
	}

	return &dto.CompareResponse{
		Period1: dto.ComparePeriodData{
			Label:        query.Value1,
			TotalIncome:  summary1.TotalIncome,
			TotalExpense: summary1.TotalExpense,
			Net:          summary1.Net,
			TopExpense:   p1TopExpense,
			TopIncome:    p1TopIncome,
		},
		Period2: dto.ComparePeriodData{
			Label:        query.Value2,
			TotalIncome:  summary2.TotalIncome,
			TotalExpense: summary2.TotalExpense,
			Net:          summary2.Net,
			TopExpense:   p2TopExpense,
			TopIncome:    p2TopIncome,
		},
		IncomeChange:  incomeChange,
		ExpenseChange: expenseChange,
		NetChange:     summary2.Net - summary1.Net,
	}, nil
}

func (s *statsService) buildCategoryNames(ctx context.Context, userID uuid.UUID) map[uuid.UUID]string {
	names := map[uuid.UUID]string{}
	if s.catRepo == nil {
		return names
	}
	categories, _, err := s.catRepo.ListByUserID(ctx, userID)
	if err != nil {
		return names
	}
	for _, c := range categories {
		names[c.ID] = c.Name
	}
	return names
}

func (s *statsService) buildSubcategoryNames(ctx context.Context, userID uuid.UUID) map[uuid.UUID]string {
	names := map[uuid.UUID]string{}
	if s.catRepo == nil || s.subRepo == nil {
		return names
	}
	categories, _, err := s.catRepo.ListByUserID(ctx, userID)
	if err != nil {
		return names
	}
	for _, c := range categories {
		subs, _, subErr := s.subRepo.ListByCategoryID(ctx, userID, c.ID)
		if subErr != nil {
			continue
		}
		for _, sub := range subs {
			names[sub.ID] = sub.Name
		}
	}
	return names
}

// Monthly returns income/expense/net per calendar month for the last N months.
func (s *statsService) Monthly(ctx context.Context, userID uuid.UUID, query dto.MonthlyQuery) (*dto.MonthlyReportResponse, error) {
	months := query.Months
	if months <= 0 || months > 60 {
		months = 12
	}

	now := time.Now().UTC()
	startMonth := time.Date(now.Year(), now.Month()-time.Month(months-1), 1, 0, 0, 0, 0, time.UTC)
	endMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)

	filter := repository.TransactionListFilter{
		Type: "all", SortBy: "transaction_date", SortDir: "asc",
		Page: 1, PerPage: 50000,
		DateFrom: &startMonth, DateTo: &endMonth,
	}
	if query.AccountID != nil && *query.AccountID != "" && *query.AccountID != "all" {
		parsedID, err := uuid.Parse(*query.AccountID)
		if err != nil {
			return nil, fmt.Errorf("statsService.Monthly.ParseAccountID: %w", err)
		}
		filter.AccountID = &parsedID
	}

	items, _, err := s.txRepo.ListByUserID(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("statsService.Monthly.List: %w", err)
	}

	categoryNames := s.buildCategoryNames(ctx, userID)
	subcategoryNames := s.buildSubcategoryNames(ctx, userID)

	type monthKey struct{ year int; month time.Month }
	type breakdownAccum map[statsBreakdownKey]float64

	monthIncome := map[monthKey]float64{}
	monthExpense := map[monthKey]float64{}
	monthExpenseBreakdown := map[monthKey]breakdownAccum{}
	monthIncomeBreakdown := map[monthKey]breakdownAccum{}

	for _, item := range items {
		mk := monthKey{year: item.TransactionDate.Year(), month: item.TransactionDate.Month()}

		catName := categoryNames[item.CategoryID]
		if catName == "" {
			catName = "Uncategorized"
		}
		subName := subcategoryNames[item.SubcategoryID]
		if subName == "" {
			subName = catName
		}
		bk := statsBreakdownKey{category: catName, subcategory: subName}

		if item.Amount >= 0 {
			monthIncome[mk] += item.Amount
			if monthIncomeBreakdown[mk] == nil {
				monthIncomeBreakdown[mk] = breakdownAccum{}
			}
			monthIncomeBreakdown[mk][bk] += item.Amount
		} else {
			monthExpense[mk] += -item.Amount
			if monthExpenseBreakdown[mk] == nil {
				monthExpenseBreakdown[mk] = breakdownAccum{}
			}
			monthExpenseBreakdown[mk][bk] += -item.Amount
		}
	}

	result := &dto.MonthlyReportResponse{Months: make([]dto.MonthlyDataPoint, 0, months)}
	for i := 0; i < months; i++ {
		t := time.Date(now.Year(), now.Month()-time.Month(months-1-i), 1, 0, 0, 0, 0, time.UTC)
		mk := monthKey{year: t.Year(), month: t.Month()}
		inc := monthIncome[mk]
		exp := monthExpense[mk]
		topN := 5
		expBreak := buildBreakdownItems(monthExpenseBreakdown[mk], exp)
		incBreak := buildBreakdownItems(monthIncomeBreakdown[mk], inc)
		if len(expBreak) > topN {
			expBreak = expBreak[:topN]
		}
		if len(incBreak) > topN {
			incBreak = incBreak[:topN]
		}
		result.Months = append(result.Months, dto.MonthlyDataPoint{
			Month:        fmt.Sprintf("%04d-%02d", t.Year(), int(t.Month())),
			TotalIncome:  inc,
			TotalExpense: exp,
			Net:          inc - exp,
			TopExpense:   expBreak,
			TopIncome:    incBreak,
		})
	}

	return result, nil
}
