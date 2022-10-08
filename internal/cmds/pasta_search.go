package cmds

import (
	"abobus/internal/services"
	"bufio"
	"math/rand"
	"os"

	twitch "github.com/gempir/go-twitch-irc/v3"
)

type PastaSearchCmd struct{}

func (PastaSearchCmd) Command() string {
	return "?pasta"
}

func (PastaSearchCmd) Description() string {
	return "Search for copypasta"
}

// TODO: cmd to add pastas
func (cmd PastaSearchCmd) Run(s *services.Services, perms []string, message twitch.PrivateMessage) (string, error) {
	// // TODO: check empty query
	// query := strings.TrimPrefix(message.Message, cmd.Command()+" ")

	// TODO: move pastes to database
	file, err := os.Open("pastes.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	pastes := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pastes = append(pastes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// TODO: fix searching
	res := pastes[rand.Intn(len(pastes))]
	return res, nil

	// type tmp struct {
	// 	Target   string
	// 	Distance int
	// }
	// matches := lo.Filter(lo.Map(pastes, func(pasta string, _ int) tmp {
	// 	return tmp{
	// 		Target:   pasta,
	// 		Distance: fuzzy.LevenshteinDistance(strings.ToLower(query), strings.ToLower(pasta)),
	// 	}
	// }), func(match tmp, _ int) bool {
	// 	return match.Distance > 0
	// })
	// if len(matches) == 0 {
	// 	return "No pastes found", nil
	// }

	// res := matches[rand.Intn(len(matches))]

	// return fmt.Sprintf("%d: %s", res.Distance, res.Target), nil
}
