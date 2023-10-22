package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FormatDate(template string, date time.Time, usePad bool) string {
	pad := func(n int) string {
		if usePad {
			return fmt.Sprintf("%02d", n)
		}
		return fmt.Sprintf("%d", n)
	}

	hour := date.Hour()

	if hour > 12 {
		hour -= 12
	}

	return strings.NewReplacer(
		"{y}", pad(date.Year()),
		"{m}", pad(int(date.Month())),
		"{d}", pad(date.Day()),
		"{h}", pad(hour),
		"{i}", pad(date.Minute()),
		"{s}", pad(date.Second()),
		"{ms}", pad(date.Nanosecond()/int(time.Millisecond)),
		"{zs}", "UTC",
		"{zd}", "+00:00",
		"{iso}", date.Format(time.RFC3339),
		"{LT}", date.Format("15:04:05"),
	).Replace(template)
}

func addDays(t time.Time, days int) time.Time {
	return t.Add(time.Duration(days) * 24 * time.Hour)
}

// Month, Day, Year
func BirthdayObject(addOne, addPad bool) (string, string, string) {
	date := GetEstDate()

	if addOne {
		date = addDays(date, 1)
	}

	m := date.Month()
	d := date.Day()
	y := RndRangeInt(1950, 2000)

	if addPad {
		return fmt.Sprintf("%02d", int(m)), fmt.Sprintf("%02d", d), strconv.Itoa(y)
	} else {
		return fmt.Sprintf("%d", int(m)), fmt.Sprintf("%d", d), strconv.Itoa(y)
	}

}

func GetEstDate() time.Time {
	est, _ := time.LoadLocation("America/New_York")
	return time.Now().In(est)
}

func Birthday(format string, addOne bool, pad bool) string {
	m, d, y := BirthdayObject(addOne, pad)
	return strings.NewReplacer(
		"{M}", m,
		"{D}", d,
		"{Y}", y,
	).Replace(format)
}
