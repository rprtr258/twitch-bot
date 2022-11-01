@_help:
  just --list --unsorted

# run locally
run:
  doppler run -- go run cmd/main.go serve

# run locally in docker
run-docker:
  doppler run --token="$(doppler configs tokens create docker --max-age 1m --plain)" -- docker compose up --build

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
  scp -i ~/.ssh/test_vds twitch-bot root@176.126.113.161:/root/go
  scp -i ~/.ssh/test_vds twitch-bot.service root@176.126.113.161:/etc/systemd/system/
  ssh vps-root 'systemctl daemon-reload'
  scp -r -i ~/.ssh/test_vds permissions.json root@176.126.113.161:/root/go/
