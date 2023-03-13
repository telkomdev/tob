package util

import (
	"fmt"
	"time"
)

// Define constant YYMMDD as the date format used in the function
const (
	YYMMDD = "2006/01/02 15:04:05"
)

// TimeDifference calculates the difference in minutes between two timestamps
// timeFrom and timeNow, passed as strings in the format defined in YYMMDD
// The function returns a string with the difference in minutes
func TimeDifference(timeFrom string, timeNow string) string {
	parsedTimeFrom, err := time.Parse(YYMMDD, timeFrom)
	if err != nil {
		return fmt.Sprintf("%d minutes", 0)
	}

	parsedTimeNow, err := time.Parse(YYMMDD, timeNow)
	if err != nil {
		return fmt.Sprintf("%d minutes", 0)
	}

	diff := parsedTimeNow.Sub(parsedTimeFrom).Minutes()
	return fmt.Sprintf("%d minutes", uint(diff))
}
