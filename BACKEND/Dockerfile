
FROM golang:1.27.1-alpine AS build


WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN mkdir -p /out && \
    CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd && \
    CGO_ENABLED=0 GOOS=linux go build -o /out/migrate ./cmd/migrate


FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=build /out/api ./api
COPY --from=build /out/migrate ./migrate

RUN mkdir uploads && chown app:app uploads

USER app

EXPOSE 8080

ENTRYPOINT ["./api"]
