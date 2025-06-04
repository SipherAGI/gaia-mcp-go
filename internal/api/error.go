package api

import (
	"fmt"
	"gaia-mcp-go/pkg/httpclient"
	"gaia-mcp-go/pkg/shared"
	"strings"
)

// ErrorKeyWord represents the error key words
type ErrorKeyWord string

const (
	ErrorKeyWordSubscriptionEnded ErrorKeyWord = "your subscription has ended"
	ErrorKeyWordCreditsExhausted  ErrorKeyWord = "no available credits"
)

// ErrorResponseMap maps the error key words to the error messages
var ErrorResponseMap = map[ErrorKeyWord]string{
	ErrorKeyWordSubscriptionEnded: fmt.Sprintf(
		"Your subscription has ended. Please update to access features here: %s/settings/account?tab=Plans&plan=subscription",
		shared.HOMEPAGE_URL,
	),
	ErrorKeyWordCreditsExhausted: fmt.Sprintf(
		"Your GAIA CREDITS balance is not enough to generate images for this task. Please buy more CREDITS here: %s/settings/account?tab=Plans&plan=credits",
		shared.HOMEPAGE_URL,
	),
}

// ProcessError processes the error and returns a new error with the appropriate message
func ProcessError(err error) error {
	if err != nil {
		// Check if the error is an API error
		if apiErr, ok := err.(*httpclient.APIError); ok {
			// Handle API errors
			msg := strings.ToLower(apiErr.Message)
			if strings.Contains(msg, string(ErrorKeyWordSubscriptionEnded)) {
				return fmt.Errorf(ErrorResponseMap[ErrorKeyWordSubscriptionEnded])
			}
			if strings.Contains(msg, string(ErrorKeyWordCreditsExhausted)) {
				return fmt.Errorf(ErrorResponseMap[ErrorKeyWordCreditsExhausted])
			}
		}

		return err
	}

	return nil
}
