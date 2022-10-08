package cmds

import (
	"abobus/internal/services"
	"database/sql"
	"fmt"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/karalef/balaboba"
	"github.com/pocketbase/dbx"
)

const _blabTable = "blab"

func responseOrLink(text, host, id string) string {
	if utf8.RuneCountInString(text) > MaxMessageLength {
		// TODO: print paste start and then link
		return fmt.Sprintf("Читать продолжение в источнике: %s", path.Join(host, "blab", id))
	}

	return text
}

type BlabGenCmd struct{}

func (BlabGenCmd) Command() string {
	return "?blab"
}

func (BlabGenCmd) Description() string {
	return "Balaboba text generation neural network"
}

func (cmd BlabGenCmd) Run(s *services.Services, perms []string, message twitch.PrivateMessage) (string, error) {
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

	id, err := s.Insert(_blabTable, map[string]any{
		"text":          response.Text,
		"author_id":     message.User.ID,
		"continuations": 1,
		"style":         styleIdx,
		"channel":       message.Channel,
	})
	if err != nil {
		return "", err
	}

	return responseOrLink(
		fmt.Sprintf("%s: %s", id, response.Text),
		s.Backend.Settings().Meta.AppUrl,
		id,
	), nil
}

type BlabContinueCmd struct{}

func (BlabContinueCmd) Command() string {
	return "?blab!"
}

func (BlabContinueCmd) Description() string {
	return "Continue paste by balaboba"
}

func (cmd BlabContinueCmd) Run(s *services.Services, perms []string, message twitch.PrivateMessage) (string, error) {
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

	response, err := s.Balaboba.Generate(text, balaboba.Style(styleIdx))
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

	return responseOrLink(
		response.Text,
		s.Backend.Settings().Meta.AppUrl,
		pasteID,
	), nil
}

type BlabReadCmd struct{}

func (BlabReadCmd) Command() string {
	return "?blab?"
}

func (BlabReadCmd) Description() string {
	return "Read pasta generated, providing it's ID"
}

func (cmd BlabReadCmd) Run(s *services.Services, perms []string, message twitch.PrivateMessage) (string, error) {
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

	return responseOrLink(
		text,
		s.Backend.Settings().Meta.AppUrl,
		pasteID,
	), nil
}
