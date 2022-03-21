package voting

import (
	"time"
)

func TimeAbsoluteDifferenceMillis(a time.Time, b time.Time) int64 {
	aMillis := a.UnixMilli()
	bMillis := b.UnixMilli()

	if aMillis >= bMillis {
		return aMillis - bMillis
	}

	return bMillis - aMillis
}

// Finds the nearest date (compared to the current time -- so yesterday, today or tomorrow) for the time (hours, minutes and seconds) of the `Time` object.
// For example, if searching for the nearest date for 00:01:00, and the time.Now() gives 23:59:59, this function will return the date for tomorrow 00:01:00.
// If time.Now() gives 00:00:01, this function will return the date for today 00:01:00.
// The day is determined by choosing the day that has the least amount difference in milliseconds compared to the current time.
func FindNearestDateForTime(Time *time.Time) *time.Time {
	hr, min, sec := Time.Clock()
	actualToday := time.Now()
	dateToday := time.Date(actualToday.Year(), actualToday.Month(), actualToday.Day(), hr, min, sec, 0, actualToday.Location())

	var dates = []time.Time{
		dateToday,                   // Today
		dateToday.AddDate(0, 0, -1), // Yesterday
		dateToday.AddDate(0, 0, 1),  // Tomorrow
	}

	nearestMillis := TimeAbsoluteDifferenceMillis(dates[0], actualToday)
	nearestDate := dates[0]
	for _, date := range dates {
		millis := TimeAbsoluteDifferenceMillis(date, actualToday)
		if millis < nearestMillis {
			nearestMillis = millis
			nearestDate = date
		}
	}

	return &nearestDate
}

func TimeStringToDate(timeString string) (*time.Time, error) {
	date, err := time.Parse(time.Kitchen, timeString)
	if err != nil {
		return nil, err
	}

	return &date, nil
}
