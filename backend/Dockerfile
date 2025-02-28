# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install necessary runtime dependencies
RUN apk update && apk add --no-cache \
    ca-certificates \
    nodejs \
    npm \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    font-noto-emoji

# Set environment variable for Chromium to run in container
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser
ENV CHROME_DISABLE_GPU 1

# Install Marp CLI
RUN npm install -g @marp-team/marp-cli

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Only copy the themes directory which contains static files
COPY --from=builder /app/services/slides/themes ./services/slides/themes

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"] 