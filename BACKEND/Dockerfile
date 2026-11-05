
FROM golang:1.27.1-alpine AS build


WORKDIR /src

# Copy Go dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the backend source code
COPY . .

RUN mkdir -p /out && \
    CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd && \
    CGO_ENABLED=0 GOOS=linux go build -o /out/migrate ./cmd/migrate


# Stage 2: Create the final lightweight image
FROM alpine:3.22

# Create a non-root user
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy the API and migration executables
COPY --from=build /out/api ./api
COPY --from=build /out/migrate ./migrate

# Create uploads directory and set permissions
RUN mkdir uploads && chown app:app uploads

# Run as a non-root user
USER app

# Backend port
EXPOSE 8080

# Start the Go application
ENTRYPOINT ["./api"]
