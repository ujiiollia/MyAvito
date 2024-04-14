FROM golang:latest AS builder
WORKDIR /usr/local/src

#dependencies
COPY go.mod go.sum ./
RUN go mod download
#USER evs:grp
#build
ADD . . 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CONFIG_PATH=./config/config.yaml go build -o ./bin/app ./cmd/app/main.go
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CONFIG_PATH=./config/config.yaml go build -o /app ./cmd/app/main.go

    
#strat
FROM alpine:latest AS runner
RUN addgroup -S grp && adduser -S evs -G grp
#USER evs
COPY --from=builder /usr/local/src/bin/app /
COPY config/config.yaml /config.yaml
#RUN groupadd -r grp && useradd -r -g evs grp

EXPOSE 8080 40000
CMD ["/bin/sh"]
#CMD ["/app"]
# CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]