package cmds

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/karalef/balaboba"
	"github.com/pocketbase/dbx"

	"abobus/internal/permissions"
	"abobus/internal/services"
)

const _blabTable = "blab"

type BlabGenCmd struct{}

func (BlabGenCmd) Command() string {
	return "?blab"
}

func (BlabGenCmd) Description() string {
	return "Balaboba text generation neural network"
}

func (cmd BlabGenCmd) Run(ctx context.Context, s *services.Services, perms permissions.PermissionsList, message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")

	if len(words) == 1 {
		return fmt.Sprintf("Available styles: %d-standart, %d-user manual, %d-recipes, %d-short stories, %d-wikipedia simplified, %d-movie synopses, %d-folk wisdom",
			balaboba.Standart,
			balaboba.UserManual,
			balaboba.Recipes,
			balaboba.ShortStories,
			balaboba.WikipediaSipmlified,
			balaboba.MovieSynopses,
			balaboba.FolkWisdom,
		), nil
	}

	// remove command
	words = words[1:]

	styleIdx, err := strconv.Atoi(words[0])
	if err != nil {
		styleIdx = int(balaboba.Standart)
	} else {
		if styleIdx < 0 || styleIdx > int(balaboba.FolkWisdom) {
			return fmt.Sprintf("Invalid style, see %s for list of available styles", cmd.Command()), nil
		}
		// remove style
		words = words[1:]
	}

	text := strings.Join(words, " ")
	response, err := s.Balaboba.Generate(text, balaboba.Style(styleIdx))
	if err != nil {
		return "", err
	}

	if response.BadQuery {
		return "PauseFish", nil
	}

	responseText := strings.ReplaceAll(response.Text, "\n", " ")

	id, err := s.Insert(_blabTable, map[string]any{
		"text":          responseText,
		"author_id":     message.User.ID,
		"continuations": 1,
		"style":         styleIdx,
		"channel":       message.Channel,
	})
	if err != nil {
		return "", err
	}

	shortResponse := fmt.Sprintf("%s %s", id, responseText)
	if utf8.RuneCountInString(shortResponse) > MaxMessageLength {
		link := path.Join(s.Backend.Settings().Meta.AppUrl, "blab", id)
		linkLen := len(link) + 1
		runes := []rune(responseText)
		upperBound := MaxMessageLength
		if utf8.RuneCountInString(responseText) < MaxMessageLength {
			upperBound = utf8.RuneCountInString(responseText)
		}
		return fmt.Sprintf("%s %s", string(runes[:upperBound-linkLen]), link), nil
	}

	return shortResponse, nil
}

type BlabContinueCmd struct{}

func (BlabContinueCmd) Command() string {
	return "?blab!"
}

func (BlabContinueCmd) Description() string {
	return "Continue paste by balaboba"
}

func (cmd BlabContinueCmd) Run(ctx context.Context, s *services.Services, perms permissions.PermissionsList, message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")

	if len(words) != 2 {
		return fmt.Sprintf("Usage: %s <paste-id>", cmd.Command()), nil
	}

	pasteID := words[1]

	db := s.Backend.DB()
	sqlCondition := dbx.And(
		dbx.NewExp("id={:id}", dbx.Params{"id": pasteID}),
		dbx.NewExp("channel={:channel}", dbx.Params{"channel": message.Channel}),
	)

	var (
		text          string
		continuations int
		styleIdx      int
	)
	if err := db.
		Select("text", "continuations", "style").
		From(_blabTable).
		Where(sqlCondition).
		Row(&text, &continuations, &styleIdx); err != nil {
		if err == sql.ErrNoRows {
			return "404: Paste not found", nil
		}

		return "", err
	}

	response, err := s.Balaboba.GenerateContext(ctx, text, balaboba.Style(styleIdx))
	if err != nil {
		return "", err
	}

	if _, err := db.
		Update(_blabTable, map[string]any{
			"text":          response.Text,
			"continuations": continuations + 1,
		}, sqlCondition).
		Execute(); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Читать продолжение в источнике: %s",
		path.Join(s.Backend.Settings().Meta.AppUrl, "blab", pasteID),
	), nil
}

type BlabReadCmd struct{}

func (BlabReadCmd) Command() string {
	return "?blab?"
}

func (BlabReadCmd) Description() string {
	return "Read pasta generated, providing it's ID"
}

func (cmd BlabReadCmd) Run(ctx context.Context, s *services.Services, perms permissions.PermissionsList, message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")

	if len(words) != 2 {
		return fmt.Sprintf("Usage: %s <paste-id>", cmd.Command()), nil
	}

	pasteID := words[1]

	db := s.Backend.DB()

	var text string
	if err := db.
		Select("text").
		From(_blabTable).
		Where(dbx.And(
			dbx.NewExp("id={:id}", dbx.Params{"id": pasteID}),
			dbx.NewExp("channel={:channel}", dbx.Params{"channel": message.Channel}),
		)).
		Row(&text); err != nil {
		if err == sql.ErrNoRows {
			return "404: Paste not found", nil
		}

		return "", err
	}

	if utf8.RuneCountInString(text) > MaxMessageLength {
		link := path.Join(s.Backend.Settings().Meta.AppUrl, "blab", pasteID)
		linkLen := len(link) + 1
		runes := []rune(text)
		return fmt.Sprintf("%s %s", string(runes[:MaxMessageLength-linkLen]), link), nil
	}

	return text, nil
}
