# ─── Stage 1: Builder ────────────────────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# build the Cobra-based CLI (which also has `serve`)
RUN CGO_ENABLED=0 GOOS=linux go build -o todo main.go

# ─── Stage 2: Final ──────────────────────────────────────────
FROM scratch
COPY --from=builder /src/todo /todo
EXPOSE 50051 8000
ENTRYPOINT ["/todo"]
CMD ["serve"]
