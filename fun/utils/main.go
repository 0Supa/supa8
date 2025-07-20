package utils

import (
	"fmt"
	"math"
	"strings"

	"github.com/0supa/supa8/config"
)

var QuoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

var privilegedUsers = make(map[string]struct{})
var botUsers = make(map[string]struct{})

func SliceToSet(slice []string, set map[string]struct{}) {
	for _, el := range slice {
		set[el] = struct{}{}
	}
}

func init() {
	SliceToSet(config.Meta.PrivilegedUsers, privilegedUsers)
	SliceToSet(config.Meta.BotUsers, botUsers)
}

func IsPrivileged(userID string) bool {
	_, exists := privilegedUsers[userID]
	return exists
}

func IsBot(userID string) bool {
	_, exists := botUsers[userID]
	return exists
}

func StringPtr(str string) *string {
	return &str
}

func FormatDuration(timeInSeconds int) string {
	hours := int(math.Floor(float64(timeInSeconds) / 3600))
	minutes := (timeInSeconds % 3600) / 60
	seconds := timeInSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
