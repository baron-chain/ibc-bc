package logging

import (
	"strings"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	eventPrefix = "events"
	typePrefix  = "type="
	separator   = "="
)

// EventsToLog converts sdk.Events into a human-readable string slice.
// It efficiently processes events and their attributes while filtering
// empty values and pre-allocating memory for better performance.
func EventsToLog(events sdk.Events) []string {
	if len(events) == 0 {
		return []string{eventPrefix}
	}

	// Pre-calculate capacity to avoid reallocation
	totalSize := 1 // for "events"
	for _, e := range events {
		totalSize++ // for type
		for _, attr := range e.Attributes {
			if len(attr.Value) > 0 {
				totalSize++
			}
		}
	}

	// Pre-allocate slice with calculated capacity
	logArgs := make([]string, 0, totalSize)
	logArgs = append(logArgs, eventPrefix)

	// Build string using StringBuilder for better performance
	var sb strings.Builder
	sb.Grow(64) // Reasonable initial size for most event strings

	for _, event := range events {
		// Add event type
		sb.WriteString(typePrefix)
		sb.WriteString(event.Type)
		logArgs = append(logArgs, sb.String())
		sb.Reset()

		// Process attributes
		for _, attr := range event.Attributes {
			if len(attr.Value) == 0 {
				continue
			}

			sb.WriteString(string(attr.Key))
			sb.WriteString(separator)
			sb.WriteString(string(attr.Value))
			logArgs = append(logArgs, sb.String())
			sb.Reset()
		}
	}

	return logArgs
}

// FormatEvent formats a single event with its attributes into a string slice.
// Useful when processing events individually.
func FormatEvent(event sdk.Event) []string {
	if len(event.Attributes) == 0 {
		return []string{typePrefix + event.Type}
	}

	logArgs := make([]string, 0, len(event.Attributes)+1)
	logArgs = append(logArgs, typePrefix+event.Type)

	var sb strings.Builder
	sb.Grow(64)

	for _, attr := range event.Attributes {
		if len(attr.Value) == 0 {
			continue
		}

		sb.WriteString(string(attr.Key))
		sb.WriteString(separator)
		sb.WriteString(string(attr.Value))
		logArgs = append(logArgs, sb.String())
		sb.Reset()
	}

	return logArgs
}

// BatchEventsToLog processes multiple event sets and combines them into a single log slice.
// Useful for aggregating logs from multiple sources.
func BatchEventsToLog(eventSets ...sdk.Events) []string {
	if len(eventSets) == 0 {
		return []string{eventPrefix}
	}

	// Calculate total capacity needed
	totalSize := 1 // for "events"
	for _, events := range eventSets {
		for _, e := range events {
			totalSize++ // for type
			for _, attr := range e.Attributes {
				if len(attr.Value) > 0 {
					totalSize++
				}
			}
		}
	}

	logArgs := make([]string, 0, totalSize)
	logArgs = append(logArgs, eventPrefix)

	var sb strings.Builder
	sb.Grow(64)

	for _, events := range eventSets {
		for _, event := range events {
			// Add event type
			sb.WriteString(typePrefix)
			sb.WriteString(event.Type)
			logArgs = append(logArgs, sb.String())
			sb.Reset()

			// Process attributes
			for _, attr := range event.Attributes {
				if len(attr.Value) == 0 {
					continue
				}

				sb.WriteString(string(attr.Key))
				sb.WriteString(separator)
				sb.WriteString(string(attr.Value))
				logArgs = append(logArgs, sb.String())
				sb.Reset()
			}
		}
	}

	return logArgs
}
