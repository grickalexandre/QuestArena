# Build frontend
FROM node:22-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Build backend
FROM golang:1.24-alpine AS backend
ENV GOTOOLCHAIN=auto
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Runtime
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend /server /app/server
COPY --from=frontend /app/frontend/dist /app/static
RUN mkdir -p /app/data
ENV PORT=8080
ENV STATIC_DIR=/app/static
ENV DEV_MODE=true
ENV DATA_DIR=/app/data
EXPOSE 8080
CMD ["/app/server"]
