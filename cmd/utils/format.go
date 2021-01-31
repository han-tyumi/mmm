package utils

import "fmt"

// FormatBigFloat formats a large float using a thousands, millions, or billion suffix.
func FormatBigFloat(value float64) string {
	switch {
	case value >= 1_000_000_000:
		return fmt.Sprintf("%.1f B", value/1_000_000_000)
	case value >= 1_000_000:
		return fmt.Sprintf("%.1f M", value/1_000_000)
	case value >= 1_000:
		return fmt.Sprintf("%.1f K", value/1_000)
	}
	return fmt.Sprint(value)
}
