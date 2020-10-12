FROM golang:1.15 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN make test
RUN make build

FROM alpine:3

RUN apk upgrade -U \
 && apk --no-cache add ca-certificates ffmpeg libva-intel-driver \
 && rm -rf /var/cache/*

COPY --from=builder /app/bin ./
CMD ["./bot"]
