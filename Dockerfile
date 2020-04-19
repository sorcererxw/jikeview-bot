FROM alpine:3

RUN apk upgrade -U \
 && apk add ca-certificates ffmpeg libva-intel-driver \
 && rm -rf /var/cache/*

COPY build ./
CMD ./jikeview-bot
