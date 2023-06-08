package transaction

import (
	"testing"

	"docqube.de/bookkeeper/pkg/services/category"
	"docqube.de/bookkeeper/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_Category_MatchesTransaction(t *testing.T) {
	categoryGroceries := category.Category{
		Name: "Groceries",
		Rules: []category.CategoryRule{
			{
				Regex:        "lidl",
				MappingField: category.MappingFieldRecipient,
			},
			{
				Regex:        "aldi",
				MappingField: category.MappingFieldRecipient,
			},
		},
	}
	categoryIncome := category.Category{
		Name: "Income",
		Rules: []category.CategoryRule{
			{
				Regex:        "gehalt",
				MappingField: category.MappingFieldBookingText,
			},
			{
				Regex:        "abrechnung",
				MappingField: category.MappingFieldPurpose,
			},
		},
	}

	tests := []struct {
		name        string
		category    *category.Category
		transaction Transaction
		want        bool
		wantErr     bool
	}{
		{
			name:     "should match lidl",
			category: &categoryGroceries,
			transaction: Transaction{
				Recipient:   utils.NewString("VISA LIDL DIENSTLEISTUNG"),
				BookingText: "Lastschrift",
				Purpose:     utils.NewString("NR XXXX 1337 MONSCHAU KAUFUMSATZ 01.01 0815123 ARN1234567890123456789"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name:     "should not match groceries",
			category: &categoryGroceries,
			transaction: Transaction{
				Recipient:   utils.NewString("VISA APPLE.COM/BILL"),
				BookingText: "Lastschrift",
				Purpose:     utils.NewString("NR XXXX 1337 KAUFUMSATZ 01.01 0815123 ARN1234567890123456789"),
			},
			want:    false,
			wantErr: false,
		},
		{
			name:     "should match booking text",
			category: &categoryIncome,
			transaction: Transaction{
				Recipient:   utils.NewString("ACME AG"),
				BookingText: "Gehalt/Rente",
				Purpose:     utils.NewString("2023/01"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name:     "should match purpose",
			category: &categoryIncome,
			transaction: Transaction{
				Recipient:   utils.NewString("ACME AG"),
				BookingText: "Gutschrift",
				Purpose:     utils.NewString("Abrechnung 2023/01"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name:     "should not match income",
			category: &categoryIncome,
			transaction: Transaction{
				Recipient:   utils.NewString("ACME AG"),
				BookingText: "Gutschrift",
				Purpose:     utils.NewString("Not matching"),
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.transaction.MatchesCategory(tt.category)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
