package internal

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/pocketbase/dbx"
)

const (
	fedCmd = "?fed"
)

func (s *Services) fed(perms []string, message twitch.PrivateMessage) (string, error) {
	words := strings.Split(message.Message, " ")
	if len(words) != 2 && len(words) != 3 {
		return fmt.Sprintf(`Usage: "%[1]s <word>" or "%[1]s <user> <word>" or "%[1]s * <word>"`, fedCmd), nil
	}

	db := s.Backend.DB()

	if len(words) == 3 && words[1] == "*" {
		theWord := words[2]
		var count int
		err := db.Select("COUNT(*) AS count").From("messages").Where(dbx.And(
			dbx.NewExp("channel={:channel}", dbx.Params{"channel": message.Channel}),
			// TODO: fix theWord = '%%' escaping
			dbx.Like("message", theWord),
			dbx.NotLike("message", fedCmd).Match(false, true),
		)).Row(&count)
		if err != nil && err != sql.ErrNoRows {
			// TODO: save log
			log.Println(err.Error())
			return "", err
		}

		return fmt.Sprintf("Слово %s было использовано %d раз.", theWord, count), nil
	}

	var (
		theWord string
		user    string
		mention string
	)
	if len(words) == 2 {
		user = message.User.Name
		theWord = words[1]
		mention = "Ты"
	} else {
		user = words[1]
		theWord = words[2]
		mention = user
	}

	var count int
	err := db.Select("COUNT(*) AS count").From("messages").Where(dbx.And(
		dbx.NewExp("channel={:channel}", dbx.Params{"channel": message.Channel}),
		// TODO: fix theWord = '%%' escaping
		dbx.Like("message", theWord),
		dbx.NotLike("message", fedCmd).Match(false, true),
		dbx.NewExp("user_name={:user_name}", dbx.Params{"user_name": user}),
	)).Row(&count)
	if err != nil && err != sql.ErrNoRows {
		// TODO: save log
		log.Println(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s использовал слово %s %d раз.", mention, theWord, count), nil
}
