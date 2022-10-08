package cmds

import (
	"context"

	twitch "github.com/gempir/go-twitch-irc/v3"

	"abobus/internal/services"
)

const MaxMessageLength = 499

type Command interface {
	// Command - get command name
	Command() string
	// Description - get command description
	Description() string
	// Run - execute command
	Run(context.Context, *services.Services, []string, twitch.PrivateMessage) (string, error)
}
