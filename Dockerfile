FROM golang:1.14 as builder


FROM alpine

RUN apk upgrade -U \
 && apk --no-cache add ca-certificates ffmpeg libva-intel-driver \
 && rm -rf /var/cache/*

COPY ./bin ./
ENTRYPOINT ["sh", "-c"]
CMD ["./bot"]
