FROM golang:1.19.2 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o twitch-bot cmd/main.go

FROM alpine:3.16.2
WORKDIR /app
COPY permissions.json pastes.txt ./
COPY --from=build /app/twitch-bot twitch-bot
EXPOSE 80
# RUN addgroup -S nonroot && adduser -S nonroot -G nonroot
# USER nonroot:nonroot

CMD ["./twitch-bot", "serve", "--http=0.0.0.0:80"]
