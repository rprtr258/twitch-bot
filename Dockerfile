FROM oven/bun:1
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY main.ts ./
COPY internal/ ./internal/
COPY permissions.json ./
EXPOSE 3000

CMD ["bun", "main.ts"]
