# Multi-stage Dockerfile
# - assets stage builds Tailwind CSS into web/static/dist.css
# - builder stage builds the Go binary
# - final stage runs the minimal Alpine image with templates and static assets

# Build frontend assets
FROM node:18-alpine AS assets
WORKDIR /app
COPY package.json package-lock.json* ./
COPY web/tailwind.css ./web/
# install dev deps and build CSS
RUN npm ci --silent
RUN npm run build:css

# Build Go binary
FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct
RUN go mod download
# copy project files and bring in built assets from previous stage
COPY . .
COPY --from=assets /app/web/static/dist.css ./web/static/dist.css
# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /app/server ./cmd/server

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
# Copy binary, templates and static assets
COPY --from=builder /app/server /app/server
COPY --from=builder /src/web/static ./web/static
COPY --from=builder /src/web/templates ./web/templates

# non-root user
RUN addgroup -S app && adduser -S -G app app
USER app
EXPOSE 8080
ENV PORT=8080
CMD ["/app/server"]
