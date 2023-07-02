package interval

import (
	"testing"
	"time"

	"docqube.de/bookkeeper/pkg/services/transaction"
)

func Test_GetFiscalMonth(t *testing.T) {
	type args struct {
		month        int
		year         int
		transactions []transaction.Transaction
	}
	tests := []struct {
		name      string
		args      args
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name: "both incomes at the beginning of the month",
			args: args{
				month: 2,
				year:  2020,
				transactions: []transaction.Transaction{
					{
						BookingDate: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 2, 4, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 3, 2, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 3, 3, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
				},
			},
			wantStart: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "both incomes at the end of the month",
			args: args{
				month: 2,
				year:  2020,
				transactions: []transaction.Transaction{
					{
						BookingDate: time.Date(2020, 1, 25, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 1, 28, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 2, 25, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 2, 26, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 3, 27, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 3, 29, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
				},
			},
			wantStart: time.Date(2020, 1, 25, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2020, 2, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "one income at the end of the month, one at the beginning",
			args: args{
				month: 2,
				year:  2020,
				transactions: []transaction.Transaction{
					{
						BookingDate: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 1, 29, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 2, 27, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 3, 28, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
				},
			},
			wantStart: time.Date(2020, 1, 29, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2020, 2, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "income at the beginning of the month, but no income in the next month",
			args: args{
				month: 2,
				year:  2020,
				transactions: []transaction.Transaction{
					{
						BookingDate: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 2, 4, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
				},
			},
			wantStart: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "both incomes at the end of the month, but no income in the next month",
			args: args{
				month: 2,
				year:  2020,
				transactions: []transaction.Transaction{
					{
						BookingDate: time.Date(2020, 1, 25, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 1, 28, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 2, 25, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 2, 26, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
				},
			},
			wantStart: time.Date(2020, 1, 25, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2020, 2, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "one income at the end of the month, one at the beginning, but no income in the next month",
			args: args{
				month: 2,
				year:  2020,
				transactions: []transaction.Transaction{
					{
						BookingDate: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 1, 29, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 2, 27, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
				},
			},
			wantStart: time.Date(2020, 1, 29, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2020, 2, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "three incomes, one at the beginning, two at the end",
			args: args{
				month: 4,
				year:  2020,
				transactions: []transaction.Transaction{
					{
						BookingDate: time.Date(2020, 3, 8, 0, 0, 0, 0, time.UTC),
						Amount:      1000.00,
					},
					{
						BookingDate: time.Date(2020, 3, 29, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 3, 31, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 4, 6, 0, 0, 0, 0, time.UTC),
						Amount:      1000.00,
					},
					{
						BookingDate: time.Date(2020, 4, 26, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
					{
						BookingDate: time.Date(2020, 4, 28, 0, 0, 0, 0, time.UTC),
						Amount:      2800.69,
					},
					{
						BookingDate: time.Date(2020, 5, 11, 0, 0, 0, 0, time.UTC),
						Amount:      1000.00,
					},
					{
						BookingDate: time.Date(2020, 5, 26, 0, 0, 0, 0, time.UTC),
						Amount:      3300.42,
					},
				},
			},
			wantStart: time.Date(2020, 3, 29, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2020, 4, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{}
			start, end := s.GetFiscalMonth(tt.args.month, tt.args.year, tt.args.transactions)

			if tt.wantStart != start {
				t.Errorf("start: got = %v, want %v", start, tt.wantStart)
			}
			if tt.wantEnd != end {
				t.Errorf("end: got = %v, want %v", end, tt.wantEnd)
			}
		})
	}
}
