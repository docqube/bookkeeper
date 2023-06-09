package category

import (
	"fmt"
	"regexp"
)

type Category struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Color       *string        `json:"color"`
	Rules       []CategoryRule `json:"rules,omitempty"`
}

type CategoryRule struct {
	ID           int64        `json:"id"`
	CategoryID   int64        `json:"categoryID"`
	MappingField MappingField `json:"mappingField"`
	Regex        string       `json:"regex"`
	Description  *string      `json:"description"`
}

type MappingField string

const (
	MappingFieldRecipient   MappingField = "recipient"
	MappingFieldBookingText MappingField = "booking_text"
	MappingFieldPurpose     MappingField = "purpose"
)

func (r *CategoryRule) Match(value string) (bool, error) {
	regex, err := regexp.Compile(fmt.Sprintf("(?i)%s", r.Regex))
	if err != nil {
		return false, err
	}
	return regex.MatchString(value), nil
}
