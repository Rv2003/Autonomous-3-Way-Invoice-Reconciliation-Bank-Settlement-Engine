
# Stage 1: Build the Go application
FROM golang:1.27.1-bookworm AS build

WORKDIR /src

# Copy Go dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the backend source code
COPY . .

# Build the Go executable into an output directory and verify it exists.
# Inspect the source files copied into the Docker build stage
RUN echo "Checking source files..." && \
    pwd && \
    find /src -maxdepth 4 -type f

# Check the Go package
RUN go list -f 'Package: {{.Name}} | Import path: {{.ImportPath}}' ./cmd

# Compile the application to a temporary executable
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /tmp/backend-api ./cmd

# Verify the executable
RUN ls -lah /tmp/backend-api && test -s /tmp/backend-api


# Stage 2: Create the final lightweight image
FROM alpine:3.22

# Create a non-root user
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy the compiled executable from Stage 1
COPY --from=build /tmp/backend-api ./api

# Create uploads directory and set permissions
RUN mkdir uploads && chown app:app uploads

# Run as a non-root user
USER app

# Backend port
EXPOSE 8080

# Start the Go application
ENTRYPOINT ["./api"]
