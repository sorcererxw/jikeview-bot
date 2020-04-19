FROM golang:1.14-alpine
COPY build ./
CMD ./jikeview-bot
