FROM golang:1.15.6-alpine3.12 AS builder

WORKDIR /go/src/github.com/conradludgate/spotify-api
ADD go.* ./
RUN go mod download
ADD *.go ./
ADD server server
RUN ls && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build ./server

FROM alpine:3.12

WORKDIR /home/spotify-proxy
ADD client client
COPY --from=builder /go/src/github.com/conradludgate/spotify-api/build spotify-proxy
ENTRYPOINT [ "spotify-proxy" ]
