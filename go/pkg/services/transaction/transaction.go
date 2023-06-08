package transaction

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"docqube.de/bookkeeper/pkg/services/category"
)

type Transaction struct {
	ID          int64              `json:"id"`
	BookingDate time.Time          `json:"bookingDate"`
	ValutaDate  time.Time          `json:"valutaDate"`
	Recipient   *string            `json:"recipient"`
	BookingText string             `json:"bookingText"`
	Purpose     *string            `json:"purpose"`
	Balance     float64            `json:"balance"`
	Amount      float64            `json:"amount"`
	Category    *category.Category `json:"category"`
}

type TransactionRequest struct {
	BookingDate *time.Time `json:"bookingDate"`
	ValutaDate  *time.Time `json:"valutaDate"`
	Recipient   *string    `json:"recipient"`
	BookingText *string    `json:"bookingText"`
	Purpose     *string    `json:"purpose"`
	Balance     *float64   `json:"balance"`
	Amount      *float64   `json:"amount"`
	CategoryID  *int64     `json:"categoryID"`
}

func (t *Transaction) Hash() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func (t *Transaction) MatchesCategory(c *category.Category) (bool, error) {
	for _, rule := range c.Rules {
		var (
			matches bool
			err     error
		)

		switch rule.MappingField {
		case category.MappingFieldRecipient:
			if t.Recipient != nil {
				matches, err = rule.Match(*t.Recipient)
			}

		case category.MappingFieldBookingText:
			matches, err = rule.Match(t.BookingText)

		case category.MappingFieldPurpose:
			if t.Purpose != nil {
				matches, err = rule.Match(*t.Purpose)
			}
		}

		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}

	return false, nil
}
