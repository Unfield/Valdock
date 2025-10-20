FROM golang:1.25.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o valdock-backend ./cmd/server/main.go

FROM alpine:3.19

RUN adduser -D -g '' valdock

WORKDIR /app
COPY --from=builder /app/valdock-backend /app/valdock-backend
RUN chown valdock:valdock /app/valdock-backend
USER valdock

ENTRYPOINT ["/app/valdock-backend"]
