FROM golang:1.17-alpine

RUN adduser -u 1000 -h /shitposta -D shitposta

WORKDIR /shitposta

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN apk add --no-cache build-base ffmpeg && \
    go build -v -o /shitposta/shitposta ./main.go

EXPOSE 80

RUN chown -R shitposta:shitposta /shitposta

USER shitposta
ENTRYPOINT ["./shitposta"]
