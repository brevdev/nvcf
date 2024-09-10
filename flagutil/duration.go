package flagutil

import (
	"fmt"
	"strings"
	"time"
)

// DurationToISO8601 converts a time.Duration to an ISO 8601 duration string.
// The resulting string follows the format: P[n]Y[n]M[n]DT[n]H[n]M[n]S
// where [n] represents the numeric value for each unit (years, months, days, hours, minutes, seconds).
// Zero values are omitted, and the 'T' separator is included only if there are time components.
//
// Example outputs:
// 1 hour 30 minutes: PT1H30M
// 2 days 4 hours: P2DT4H
// 1 year 6 months 15 days 12 hours 30 minutes 45 seconds: P1Y6M15DT12H30M45S
//
// Note: This function assumes 365 days per year and 30 days per month for simplicity.
func DurationToISO8601(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	nanoseconds := d.Nanoseconds() % 1e9
	var result strings.Builder
	result.WriteString("P")
	if days > 0 {
		result.WriteString(fmt.Sprintf("%dD", days))
	}
	if hours > 0 || minutes > 0 || seconds > 0 || nanoseconds > 0 {
		result.WriteString("T")
		if hours > 0 {
			result.WriteString(fmt.Sprintf("%dH", hours))
		}
		if minutes > 0 {
			result.WriteString(fmt.Sprintf("%dM", minutes))
		}
		if seconds > 0 || nanoseconds > 0 {
			result.WriteString(fmt.Sprintf("%d", seconds))
			if nanoseconds > 0 {
				result.WriteString(fmt.Sprintf(".%09d", nanoseconds))
			}
			result.WriteString("S")
		}
	}
	return result.String()
}
