FROM golang:1.18.3 as build

COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o minipolly

FROM alpine:3.6 as runtime
RUN apk add --no-cache ffmpeg sox soxi

COPY --from=build /app/minipolly /app/minipolly

CMD ["/app/minipolly"]