# ─── Stage 1: Builder ────────────────────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /todo main.go

# ─── Stage 2: Runtime ───────────────────────────────────────
FROM alpine:3.18 AS runtime

# install a tiny shell + utils so we can keep container alive
RUN apk add --no-cache bash coreutils
COPY --from=builder /todo /usr/local/bin/todo

# keep the container alive
ENTRYPOINT ["tail", "-f", "/dev/null"]

CMD []
