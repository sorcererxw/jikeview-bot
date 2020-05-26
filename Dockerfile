FROM golang as build

WORKDIR /app
COPY ./ ./
RUN make build

FROM alpine

RUN apk upgrade -U \
 && apk --no-cache add ca-certificates ffmpeg libva-intel-driver \
 && rm -rf /var/cache/*

COPY --from=build /app/bin ./
ENTRYPOINT ["sh", "-c"]
CMD ["./bot"]
