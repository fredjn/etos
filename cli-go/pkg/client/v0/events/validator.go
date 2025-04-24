package events

import (
	"fmt"

	"github.com/eiffel-community/eiffelevents-sdk-go"
	"github.com/eiffel-community/eiffelevents-sdk-go/validator"
)

// ValidateEvent validates an Eiffel event using the SDK's validator
func ValidateEvent(event *eiffelevents.Any) error {
	// Get the default validator set
	v := validator.DefaultSet()

	// Validate the event
	if err := v.Validate(event); err != nil {
		return fmt.Errorf("event validation failed: %w", err)
	}

	return nil
}

// ValidateEvents validates a collection of Eiffel events
func ValidateEvents(events *Events) error {
	// Validate TERCC event if present
	if events.TERCC != nil {
		if err := ValidateEvent(events.TERCC); err != nil {
			return fmt.Errorf("TERCC event validation failed: %w", err)
		}
	}

	// Validate activity events if present
	if events.Activity != nil {
		if events.Activity.Triggered != nil {
			if err := ValidateEvent(events.Activity.Triggered); err != nil {
				return fmt.Errorf("activity triggered event validation failed: %w", err)
			}
		}
		if events.Activity.Canceled != nil {
			if err := ValidateEvent(events.Activity.Canceled); err != nil {
				return fmt.Errorf("activity canceled event validation failed: %w", err)
			}
		}
		if events.Activity.Finished != nil {
			if err := ValidateEvent(events.Activity.Finished); err != nil {
				return fmt.Errorf("activity finished event validation failed: %w", err)
			}
		}
	}

	// Validate test suite events
	for i, suite := range events.MainSuites {
		if suite.Started != nil {
			if err := ValidateEvent(suite.Started); err != nil {
				return fmt.Errorf("test suite %d started event validation failed: %w", i, err)
			}
		}
		if suite.Finished != nil {
			if err := ValidateEvent(suite.Finished); err != nil {
				return fmt.Errorf("test suite %d finished event validation failed: %w", i, err)
			}
		}
	}

	return nil
} 