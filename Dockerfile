#FROM golang:alpine as test
#
#WORKDIR /app
#RUN apk upgrade -U \
# && apk --no-cache --update add make ca-certificates ffmpeg libva-intel-driver gcc \
# && rm -rf /var/cache/*
#COPY ./ ./
#RUN make test

#FROM golang as build
#
#WORKDIR /app
#COPY ./ ./
#RUN make build

FROM alpine

RUN apk upgrade -U \
 && apk --no-cache add ca-certificates ffmpeg libva-intel-driver \
 && rm -rf /var/cache/*

COPY ./bin ./
ENTRYPOINT ["sh", "-c"]
CMD ["./bot"]
