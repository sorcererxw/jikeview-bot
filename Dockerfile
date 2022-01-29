FROM alpine:3

RUN apk upgrade -U \
 && apk --no-cache add ca-certificates ffmpeg libva-intel-driver \
 && rm -rf /var/cache/*

COPY bin ./
CMD ["./bot"]
