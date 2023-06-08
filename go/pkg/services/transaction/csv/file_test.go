package csv

import (
	"io"
	"os"
	"testing"
	"time"

	"docqube.de/bookkeeper/pkg/services/transaction"
	"docqube.de/bookkeeper/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_ParseFile(t *testing.T) {
	tests := []struct {
		name    string
		reader  func() io.Reader
		config  FileConfig
		want    []transaction.Transaction
		wantErr bool
	}{
		{
			name: "should parse ing file",
			reader: func() io.Reader {
				file, err := os.Open("./testing/ing.csv")
				if err != nil {
					t.Errorf("reading test file: %s", err)
				}
				return file
			},
			config: INGConfig,
			want: []transaction.Transaction{
				{
					BookingDate: time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					ValutaDate:  time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					Recipient:   utils.NewString("VISA KAUFLAND MONSCHAU 8710"),
					BookingText: "Lastschrift",
					Purpose:     utils.NewString("NR XXXX 0815 MONSCHAU Apple Pay"),
					Balance:     500,
					Amount:      -13.37,
					Category:    nil,
				},
				{
					BookingDate: time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					ValutaDate:  time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					Recipient:   utils.NewString("Jan Muster"),
					BookingText: "Überweisung",
					Purpose:     nil,
					Balance:     513.37,
					Amount:      -29,
					Category:    nil,
				},
				{
					BookingDate: time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					ValutaDate:  time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					Recipient:   utils.NewString("Max Muster"),
					BookingText: "Gutschrift",
					Purpose:     nil,
					Balance:     542.37,
					Amount:      150,
					Category:    nil,
				},
				{
					BookingDate: time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					ValutaDate:  time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					Recipient:   utils.NewString("AUTO BANK AG NL Deutschland"),
					BookingText: "Lastschrift",
					Purpose:     utils.NewString("Auto Leasing/VT12345678 05/23 Rate"),
					Balance:     392.37,
					Amount:      -69.42,
					Category:    nil,
				},
				{
					BookingDate: time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					ValutaDate:  time.Date(2023, time.May, 22, 0, 0, 0, 0, time.UTC),
					Recipient:   utils.NewString("Telekom Deutschland GmbH"),
					BookingText: "Lastschrift",
					Purpose:     utils.NewString("Mobilfunk Kundenkonto 123456789 RG 9876543210123456789/01.01.2023"),
					Balance:     461.79,
					Amount:      -25.99,
					Category:    nil,
				},
				{
					BookingDate: time.Date(2023, time.May, 19, 0, 0, 0, 0, time.UTC),
					ValutaDate:  time.Date(2023, time.May, 19, 0, 0, 0, 0, time.UTC),
					Recipient:   utils.NewString("VISA DM-DROGERIE MARKT"),
					BookingText: "Lastschrift",
					Purpose:     utils.NewString("NR XXXX 1234 MONSCHAU KAUFUMSATZ 01.01 123456789"),
					Balance:     487.78,
					Amount:      -13.98,
					Category:    nil,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := tt.reader()
			got, err := ParseFile(reader, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))
			assert.Equal(t, tt.want, got)
		})
	}
}
