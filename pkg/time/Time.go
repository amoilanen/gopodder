package time

import "time"

func tryParseWithLayouts(timeInput string, formats ...string) (time.Time, error) {
	var parsed time.Time
	var err error
	for _, layout := range formats {
		parsed, err = time.Parse(layout, timeInput)
		if err == nil {
			return parsed, nil
		}
	}
	return parsed, err
}

func CurrentTime() string {
	currentTime := time.Now()
	return currentTime.Format(time.RFC3339)
}

func ParseTime(timeInput string) time.Time {
	parsedDate, err := tryParseWithLayouts(timeInput, time.RFC3339, time.DateOnly, time.RFC1123, time.RFC1123Z)
	if err != nil {
		panic(err)
	}
	return parsedDate
}
