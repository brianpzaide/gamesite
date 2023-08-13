FROM golang:1.19-alpine AS builder
RUN apk update && apk add make git
WORKDIR /home/gamesite
COPY . .
RUN go install github.com/knadh/stuffbin/...

RUN make build

FROM golang:1.19-alpine AS deploy
RUN apk --no-cache add ca-certificates
WORKDIR /home/gamesite
COPY --from=builder /home/gamesite/stuffedgamesite .
COPY config.env .
CMD ["./stuffedgamesite"]