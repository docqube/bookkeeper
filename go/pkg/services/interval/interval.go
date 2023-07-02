package interval

import "time"

type FiscalMonth struct {
	Month int       `json:"month"`
	Year  int       `json:"year"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
