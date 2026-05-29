# Single-image monolith: builds the React SPA and the Go API, then serves both
# from one process on port 8080. Build context must be the repository root.

# ---- Frontend build ----
FROM node:22-alpine AS frontend
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ---- Backend build ----
FROM golang:1.24-alpine AS backend
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

# ---- Runtime ----
# The binary is statically linked (CGO disabled) and only talks to PostgreSQL
# over TCP, so no extra OS packages are required. adduser is a BusyBox builtin.
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
COPY --from=backend /api /api
COPY --from=frontend /app/dist /web
# STATIC_DIR makes the Go server also serve the built SPA (single-port monolith).
ENV STATIC_DIR=/web
USER appuser
EXPOSE 8080
ENTRYPOINT ["/api"]
