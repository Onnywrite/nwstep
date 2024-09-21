FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./bin/app ./cmd/main.go

FROM alpine:3.20 AS runner
WORKDIR /bin/app

COPY --from=builder /app/bin ./
RUN adduser -DH usr && chown -R usr: /bin/app && chmod -R 700 /bin/app

USER usr
 
CMD [ "./app" ]