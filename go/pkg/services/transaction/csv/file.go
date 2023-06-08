package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"docqube.de/bookkeeper/pkg/services/transaction"
	"golang.org/x/text/transform"
)

var (
	ErrInvalidLineLength = fmt.Errorf("invalid line length")
)

func ParseFile(reader io.Reader, config FileConfig) ([]transaction.Transaction, error) {
	var r *csv.Reader

	if config.FileEncoding != nil {
		decodingReader := transform.NewReader(reader, config.FileEncoding.NewDecoder())
		r = csv.NewReader(decodingReader)
	} else {
		r = csv.NewReader(reader)
	}

	r.Comma = config.Delimiter
	r.FieldsPerRecord = config.FieldsPerRecord
	skippedHeader := false

	transactions := make([]transaction.Transaction, 0)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// skip lines with invalid length
			if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
				continue
			}
			return nil, err
		}

		// as possibly need to iterate over some lines that are neither header nor data
		// we need to skip the header here and set an extra flag and can not rely on the
		// csv package to skip the header or use the line number
		if config.HasHeader && !skippedHeader {
			skippedHeader = true
			continue
		}

		transaction, err := parseRecord(record, config)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, *transaction)
	}

	return transactions, nil
}

func parseRecord(record []string, config FileConfig) (*transaction.Transaction, error) {
	bookingDate, err := time.Parse(config.DateFormat, record[config.BookingDate])
	if err != nil {
		return nil, err
	}

	valutaDate, err := time.Parse(config.DateFormat, record[config.ValutaDate])
	if err != nil {
		return nil, err
	}

	recipient := &record[config.Recipient]
	if *recipient == "" {
		recipient = nil
	}

	bookingText := record[config.BookingText]
	if bookingText == "" {
		return nil, fmt.Errorf("booking text is empty")
	}

	purpose := &record[config.Purpose]
	if *purpose == "" {
		purpose = nil
	}

	rawBalance := record[config.Balance]
	if rawBalance == "" {
		return nil, fmt.Errorf("balance is empty")
	}
	rawBalance = strings.ReplaceAll(rawBalance, string(config.NumberFormat.ThousandSeparator), "")
	rawBalance = strings.ReplaceAll(rawBalance, string(config.NumberFormat.DecimalSeparator), ".")
	balance, err := strconv.ParseFloat(rawBalance, 64)
	if err != nil {
		return nil, err
	}

	rawAmount := record[config.Amount]
	if rawAmount == "" {
		return nil, fmt.Errorf("amount is empty")
	}
	rawAmount = strings.ReplaceAll(rawAmount, string(config.NumberFormat.ThousandSeparator), "")
	rawAmount = strings.ReplaceAll(rawAmount, string(config.NumberFormat.DecimalSeparator), ".")
	amount, err := strconv.ParseFloat(rawAmount, 64)
	if err != nil {
		return nil, err
	}

	return &transaction.Transaction{
		BookingDate: bookingDate,
		ValutaDate:  valutaDate,
		Recipient:   recipient,
		BookingText: bookingText,
		Purpose:     purpose,
		Balance:     balance,
		Amount:      amount,
	}, nil
}
