@_help:
  just --list --unsorted

# run server
run:
  rwenv -i -e .env go run cmd/main.go serve

# generates: [dataset]
get_chat_logs:
  python3 dump.py >> tmp
  cat dataset tmp | sort -u > tmp_dataset
  mv tmp_dataset dataset
  rm tmp

# sources: [dataset]
# generates: [model.json]
train: get_chat_logs
  python quotes_generator/main.py train --dataset dataset --model model.json

# bump dependencies
@bump:
  go get -u ./...
  go mod tidy

# check todos
@todo:
  rg 'TODO' --glob '**/*.go' || echo 'All done!'

# stupid deploy automation
deploy:
  GOOS=linux GOARCH=amd64 go build -o ./twitch-bot cmd/main.go
  scp -i ~/.ssh/test_vds twitch-bot root@176.126.113.161:/root/go
  scp -i ~/.ssh/test_vds twitch-bot.service root@176.126.113.161:/etc/systemd/system/
  ssh root@176.126.113.161 -i ~/.ssh/test_vds 'systemctl daemon-reload'
  scp -r -i ~/.ssh/test_vds permissions.json root@176.126.113.161:/root/go/
