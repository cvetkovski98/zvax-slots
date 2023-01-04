package delivery

import "errors"

var ErrInvalidDateFormat = errors.New("Invalid date format. Expected format: 2006-01-02")
var ErrInvalidDateRange = errors.New("Invalid date range. Start date must be before end date")
