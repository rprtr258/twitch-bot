package internal

import (
	"fmt"
	"strconv"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v3"
	"github.com/karalef/balaboba"
)

// last_balaboba = ""

// def balabob(text, skip=0):
//     from subprocess import check_output
//     TO_SKIP = len("please wait up to 15 seconds Без стиля".split())
//     output = check_output(["./balaboba"] + text.split()).decode("utf-8").split()
//     print(output)
//     if "на острые темы, например, про политику или религию" in ' '.join(output):
//         return "PauseFish"
//     last_balaboba = ' '.join(output[TO_SKIP + skip:])
//     return last_balaboba

// @app.route("/blab/<idd>")
// def long_blab(idd):
//     db = read_db()
//     db[idd] = balabob(db[idd])
//     load_db(db)
//     return f'''<p style="padding: 10% 15%; font-size: 1.8em;">{db[idd]}</p>'''

// @app.route("/b")
// def balaboba():
//     message = request.args.get("m")
//     response = balabob(message, skip=len(message.split()))
//     if len(response) < 300:
//         print("TOO SHORT: ", message, len(response))
//         return response
//     else:
//         db = read_db()
//         idd = max(map(int, db.keys())) + 1
//         db[idd] = message + ' ' + response
//         load_db(db)
//         return f"Читать продолжение в источнике: secure-waters-73337.herokuapp.com/blab/{idd}"

const (
	blabCmd = "?blab"
)

func (s *Services) blab(message twitch.PrivateMessage) (string, error) {
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
			return fmt.Sprintf("Invalid style, see %s for list of available styles", blabCmd), nil
		}
		// remove style
		words = words[1:]
	}

	text := strings.Join(words, " ")
	response, err := s.Balaboba.Generate(text, balaboba.Style(styleIdx))
	if err != nil {
		return "", err
	}

	return response.Text, nil

	// db := s.Backend.DB()

	// theWord := words[1]

	// var count int
	// err = db.Select("COUNT(*) AS count").From("messages").Where(dbx.And(
	// 	dbx.NewExp("channel={:channel}", dbx.Params{"channel": message.Channel}),
	// 	// TODO: fix theWord = '%%' escaping
	// 	dbx.Like("message", theWord),
	// 	dbx.NotLike("message", fedCmd).Match(false, true),
	// 	dbx.NewExp("user_id={:user_id}", dbx.Params{"user_id": userID}),
	// )).Row(&count)
	// if err != nil && err != sql.ErrNoRows {
	// 	// TODO: save log
	// 	log.Println(err.Error())
	// 	return "", err
	// }

	// // TODO: add feed after duration
	// return fmt.Sprintf("Ты использовал слово %s %d раз.", theWord, count), nil
}
