package events

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/eventbridge"
)

// RuleEnabledFromState infers from its state whether or not a rule is enabled.
func RuleEnabledFromState(state string) (bool, error) {
	if state == eventbridge.RuleStateEnabled {
		return true, nil
	}

	if state == eventbridge.RuleStateDisabled {
		return false, nil
	}

	// We don't just blindly trust AWS as they tend to return
	// unexpected values in similar cases (different casing etc.)
	return false, fmt.Errorf("unable to infer enabled from state: %s", state)
}

// RuleStateFromEnabled returns a rule's state based on whether or not it is enabled.
func RuleStateFromEnabled(enabled bool) string {
	if enabled {
		return eventbridge.RuleStateEnabled
	}

	return eventbridge.RuleStateDisabled
}
