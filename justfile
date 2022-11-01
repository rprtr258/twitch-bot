@_help:
  just --list --unsorted

# run locally
run:
  doppler run -- go run cmd/main.go serve

# run locally in docker
run-docker:
  doppler run --token="$(doppler configs tokens create docker --max-age 1m --plain)" -- docker compose up --build --remove-orphans

# bump dependencies
@bump:
  go get -u ./...
  go mod tidy

# check todos
@todo:
  rg 'TODO' --glob '**/*.go' || echo 'All done!'

# TODO: update rules

# sources: [dataset]
# generates: [model.json]
# train: get_chat_logs
#   python quotes_generator/main.py train --dataset dataset --model model.json

# stupid deploy automation
deploy:
  ssh vps 'cd twitch-bot-deploy && docker compose pull && doppler run -- docker compose up'
