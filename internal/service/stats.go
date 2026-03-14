package service

import (
	"context"
	"fmt"
	"ledgerA/internal/dto"
	"ledgerA/internal/repository"
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
