FROM golang:latest AS builder
WORKDIR /usr/local/src

#dependencies
COPY go.mod go.sum ./
RUN go mod download

#build
ADD . . 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CONFIG_PATH=./config/config.yaml go build -o ./bin/app ./cmd/app/main.go

    
#strat
FROM alpine:latest AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY config/config.yaml /config.yaml
CMD ["/app --config=config.yaml"]