FROM golang:1.19.2 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o twitch-bot cmd/main.go

FROM alpine:3.16.2

RUN wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler

RUN apk add python3
WORKDIR /app
COPY --from=build /app/twitch-bot twitch-bot
EXPOSE 80
# RUN addgroup -S nonroot && adduser -S nonroot -G nonroot
# USER nonroot:nonroot

# ENTRYPOINT ["./twitch-bot", "serve"]
CMD ["doppler", "run", "--", "./twitch-bot", "serve", "--http=:80"]
