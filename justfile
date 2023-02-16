@_help:
  just --list --unsorted

# run locally
run:
  rwenv -ie .env -- go run cmd/main.go serve

# run locally in docker
run-docker:
  rwenv -ie .env docker compose up --build --remove-orphans

# bump dependencies
@bump:
  go get -u ./...
  go mod tidy

# check todos
@todo:
  rg 'TODO' --glob '**/*.go' || echo 'All done!'

# sources: [dataset]
# generates: [model.json]
# train: get_chat_logs
#   python quotes_generator/main.py train --dataset dataset --model model.json
