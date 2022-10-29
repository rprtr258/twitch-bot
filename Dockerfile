FROM golang:1.19.2 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY permissions.json pastes.txt ./
COPY cmd/ internal/ ./
RUN go build -o twitch-bot

FROM golang:1.19.2-alpine

RUN apk add python3
WORKDIR /
COPY pyth/ /
COPY --from=build /twitch-bot /twitch-bot
EXPOSE 80
USER nonroot:nonroot

ENTRYPOINT ["/twitch-bot"]
