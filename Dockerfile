FROM alpine:3.16.2

WORKDIR /app

COPY permssions.json permssions.json
RUN cat permssions.json
# EXPOSE 80
# RUN addgroup -S nonroot && adduser -S nonroot -G nonroot
# USER nonroot:nonroot

# ENTRYPOINT ["./twitch-bot", "serve"]
# CMD ["doppler", "run", "--", "./twitch-bot", "serve", "--http=:80"]
