package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	Type EventType `json:"type"`
	Text string    `json:"text"`
	Date time.Time `json:"date"`
}

func main() {
	var calendar []Event
	_ = calendar

	var msgs []string = []string{"Birthday \"Joe Soap\" 9/5/1990", "Birthday \"Chicken Lick'n\" 1/22/2000", "Birthday \"Howzit Brew\" 02/13/1995", "Birthday \"Another Day\" 02/14/1995"}

	for _, msg := range msgs {
		event, err := EventFromMsg(msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		AddEvent(&calendar, event)
		fmt.Println(event)
	}

	_ = ExportEventsToFile(&calendar, "events-store")
}

func EventFromMsg(msg string) (Event, error) {
	eventDate, dateParsingErr := ParseDateFromMsg(msg)
	if dateParsingErr != nil {
		return Event{}, dateParsingErr
	}

	eventText, textParsingErr := ParseTextFromMsg(msg)
	if textParsingErr != nil {
		return Event{}, textParsingErr
	}

	eventType := ParseEventTypeFromMsg(msg)

	return Event{Type: eventType, Text: eventText, Date: eventDate}, nil
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

func ParseTextFromMsg(msg string) (string, error) {
	// Extract the 'text' component from a message
	r := regexp.MustCompile(`"([^"]*)"`)

	matches := r.FindStringSubmatch(msg)

	if matches == nil {
		return "", fmt.Errorf("No 'text' component found!")
	}

	return matches[1], nil
}

func ParseDateFromMsg(msg string) (time.Time, error) {
	// Match on DD/MM/YYYY or D/M/YYYY (or combination DD, D, MM, M)
	r := regexp.MustCompile(`(\d{1,2})/(\d{1,2})/(\d{4})`)

	matches := r.FindStringSubmatch(msg)
	if matches == nil {
		return time.Time{}, fmt.Errorf("failed to parse DD/MM/YYYY")
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

func AddEvent(calendar *[]Event, event Event) {
	*calendar = append(*calendar, event)
}

func ExportEventsToFile(calendar *[]Event, fileName string) error {
	exportFile, err := os.Create(fileName + ".json")
	if err != nil {
		return fmt.Errorf("could not create export file")
	}
	defer exportFile.Close()

	jsonEncoder := json.NewEncoder(exportFile)
	jsonEncoder.SetIndent("", "    ")

	for _, event := range *calendar {
		encodeErr := jsonEncoder.Encode(event)
		if encodeErr != nil {
			fmt.Println("encoding error for event:", event)
		}
	}
	return nil
}
