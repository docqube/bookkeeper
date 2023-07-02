package interval

import (
	"database/sql"
	"time"

	"docqube.de/bookkeeper/pkg/services/transaction"
	"docqube.de/bookkeeper/pkg/utils"
)

type Service struct {
	transactionService *transaction.Service
}

func NewService(db *sql.DB) *Service {
	return &Service{
		transactionService: transaction.NewService(db),
	}
}

func (s *Service) GetFiscalMonthWithIncomeCategoryID(month int, year int, incomeCategoryID int64) (*time.Time, *time.Time, error) {
	previousMonthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	previousMonthStart.AddDate(0, -1, 0)
	nextMonthEnd := previousMonthStart.AddDate(0, 3, -1)

	incomeTransactions, err := s.transactionService.ListByCategoryID(previousMonthStart, nextMonthEnd, incomeCategoryID, transaction.OrderByDirectionAsc)
	if err != nil {
		return nil, nil, err
	}

	start, end := s.GetFiscalMonth(month, year, incomeTransactions.Items)
	return &start, &end, nil
}

func (s *Service) GetFiscalMonth(month int, year int, incomeTransactions []transaction.Transaction) (time.Time, time.Time) {
	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}

	fiscalStart := s.CalculateFiscalStartDate(month, year, incomeTransactions)
	fiscalEnd := s.CalculateFiscalStartDate(nextMonth, nextYear, incomeTransactions)

	// return one day subtracted from the end date, as we don't want to
	// include the next month with income in the requested fiscal month
	return fiscalStart, fiscalEnd.AddDate(0, 0, -1)
}

func (s *Service) CalculateFiscalStartDate(month int, year int, incomeTransactions []transaction.Transaction) time.Time {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	distanceMap := make(map[int]transaction.Transaction)
	distanceValues := make([]int, 0)

	for _, t := range incomeTransactions {
		difference := int(t.BookingDate.Sub(startOfMonth).Hours() / 24)
		distanceMap[difference] = t
		distanceValues = append(distanceValues, difference)
	}

	for _, diff := range distanceValues {
		if diff < 0 && utils.Abs(diff) >= 15 {
			// ignore transactions that are more than 15 days
			// before the start of the month
			continue
		}
		return distanceMap[diff].BookingDate
	}

	return startOfMonth
}
