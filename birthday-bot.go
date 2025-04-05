package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type EventType int

const (
	Generic EventType = iota
	Birthday
)

type Event struct {
	Type EventType
	Text string
	Date time.Time
}

func main() {
	var calendar [366][]Event
	_ = calendar

	msg := "Birthday \"Joe Soap\" 9/5/1990"
	bDate, parsingErr := ParseDateFromMsg(msg)
	if parsingErr == nil {
		fmt.Println(bDate)
	} else {
		fmt.Println("error! ", parsingErr)
	}

	fmt.Println("Event type: ", ParseEventTypeFromMsg(msg))

}

func ParseEventTypeFromMsg(msg string) EventType {
	// Extract the type of event from a message
	birthdayRegex := regexp.MustCompile(`(?i)birthday`)

	// Check if the text contains 'birthday' (case-insensitive)
	containsBirthday := birthdayRegex.MatchString(msg)

	if containsBirthday {
		return Birthday
	}

	return Generic
}

func ParseTextFromMsg(msg string) string {
	// Extract the 'text' component from a message
	return ""
}

func ParseDateFromMsg(msg string) (time.Time, error) {
	// Match on DD/MM/YYYY or D/M/YYYY (or combination DD, D, MM, M)
	r := regexp.MustCompile(`(\d{1,2})/(\d{1,2})/(\d{4})`)

	matches := r.FindStringSubmatch(msg)
	if matches == nil {
		return time.Time{}, fmt.Errorf("failed to parse DD/MM/YYYY")
	}

	if len(matches) != 4 {
		return time.Time{}, fmt.Errorf("multiple dates found")
	}

	day, dayCnvErr := strconv.Atoi(matches[1])
	month, monthCnvErr := strconv.Atoi(matches[2])
	year, yearCnvErr := strconv.Atoi(matches[3])

	if dayCnvErr != nil || monthCnvErr != nil || yearCnvErr != nil {
		return time.Time{}, fmt.Errorf("invalid date provided!")
	}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return date, nil
}

func AddEvent(calendar *[366][]Event, event Event) {
	dayOfYear := event.Date.YearDay()
	calendar[dayOfYear-1] = append(calendar[dayOfYear-1], event)
}
