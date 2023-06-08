package transaction

import (
	"testing"

	"docqube.de/bookkeeper/pkg/services/category"
	"docqube.de/bookkeeper/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_MatchTransactionCategory(t *testing.T) {
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

	service := &Service{
		categories: []category.Category{categoryGroceries, categoryIncome},
	}

	tests := []struct {
		name         string
		service      *Service
		transaction  Transaction
		wantCategory *category.Category
		wantErr      bool
	}{
		{
			name:    "should match groceries",
			service: service,
			transaction: Transaction{
				Recipient:   utils.NewString("VISA LIDL DIENSTLEISTUNG"),
				BookingText: "Lastschrift",
				Purpose:     utils.NewString("NR XXXX 1337 MONSCHAU KAUFUMSATZ 01.01 0815123 ARN1234567890123456789"),
			},
			wantCategory: &categoryGroceries,
			wantErr:      false,
		},
		{
			name:    "should match income",
			service: service,
			transaction: Transaction{
				Recipient:   utils.NewString("ACME AG"),
				BookingText: "Gehalt/Rente",
				Purpose:     utils.NewString("2023/01"),
			},
			wantCategory: &categoryIncome,
			wantErr:      false,
		},
		{
			name:    "should not match any category",
			service: service,
			transaction: Transaction{
				Recipient:   utils.NewString("ACME AG"),
				BookingText: "Gutschrift",
				Purpose:     utils.NewString("Not matching"),
			},
			wantCategory: nil,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, err := tt.service.MatchTransactionCategory(&tt.transaction)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCategory, category)
		})
	}
}
