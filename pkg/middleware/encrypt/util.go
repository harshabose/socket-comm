package encrypt

import (
	"fmt"
	"time"
)

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	} else if d < time.Minute {
		return fmt.Sprintf("%.1f s", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1f min", d.Minutes())
	} else {
		return fmt.Sprintf("%.1f h", d.Hours())
	}
}
