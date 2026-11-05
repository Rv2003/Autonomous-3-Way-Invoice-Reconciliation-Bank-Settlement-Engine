# Stage 1: Build the Go application
FROM golang:1.27.1-bookworm AS build

WORKDIR /src

# Copy Go dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the backend source code
COPY . .

# Compile both the API and the migration runner
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /out/api ./cmd && \
    CGO_ENABLED=0 GOOS=linux go build -v -o /out/migrate ./cmd/migrate

# Stage 2: Create the final lightweight image
FROM alpine:3.22

# Create a non-root user
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy both executables from build stage
COPY --from=build /out/api ./api
COPY --from=build /out/migrate ./migrate

# Create uploads directory and set permissions
RUN mkdir -p uploads && chown -R app:app uploads

# Run as a non-root user
USER app

# Backend port
EXPOSE 8080

# Start the Go application
ENTRYPOINT ["./api"]
