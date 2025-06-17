# Stage 1: Build the Go app
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY . .



# Add this line to build a fully static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-api

# Stage 2: Run in small image
FROM alpine:latest

# Install tzdata for timezone information
RUN apk add --no-cache tzdata

WORKDIR /root/
COPY --from=builder /app/weather-api .

# Set the desired timezone
ENV TZ=America/Toronto

EXPOSE 8080
CMD ["./weather-api"]