package csv

import "golang.org/x/text/encoding/charmap"

type FileConfig struct {
	Delimiter       rune
	FileEncoding    *charmap.Charmap
	FieldsPerRecord int
	HasHeader       bool
	DateFormat      string
	NumberFormat    NumberFormat
	BookingDate     int
	ValutaDate      int
	Recipient       int
	BookingText     int
	Purpose         int
	Balance         int
	Amount          int
}

type NumberFormat struct {
	DecimalSeparator  rune
	ThousandSeparator rune
}

var INGConfig = FileConfig{
	Delimiter:       ';',
	FileEncoding:    charmap.Windows1252,
	FieldsPerRecord: 9,
	HasHeader:       true,
	DateFormat:      "02.01.2006",
	NumberFormat:    NumberFormat{DecimalSeparator: ',', ThousandSeparator: '.'},
	BookingDate:     0,
	ValutaDate:      1,
	Recipient:       2,
	BookingText:     3,
	Purpose:         4,
	Balance:         5,
	Amount:          7,
}
