# ---- Build stage ----
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

# ---- Runtime stage ----
# The binary is statically linked (CGO disabled) and only talks to PostgreSQL
# over TCP, so no extra OS packages are required. adduser is a BusyBox builtin.
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
COPY --from=build /api /api
USER appuser
EXPOSE 8080
ENTRYPOINT ["/api"]
