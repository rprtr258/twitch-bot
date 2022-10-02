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
